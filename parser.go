package flv

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Parser struct {
	closeOnce sync.Once
	stop      chan struct{}

	avcHeader []byte
}

func NewParser() *Parser {
	return &Parser{
		stop: make(chan struct{}),
	}
}

func (p *Parser) New() *Parser {
	return &Parser{
		stop: make(chan struct{}),
	}
}

func (p *Parser) Stop() {
	p.closeOnce.Do(func() {
		close(p.stop)
	})
}

func (p *Parser) Type() string {
	return "flv"
}

func (p *Parser) Parse(streamUrl string, out string) (err error) {
	req, err := http.NewRequest("GET", streamUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "Chrome/59.0.3071.115")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d, err := NewDemuxer(resp.Body)
	if err != nil {
		return err
	}
	defer d.Close()

	o, err := os.Create(out)
	if err != nil {
		return err
	}

	m, err := NewMuxer(o)
	if err != nil {
		return err
	}
	defer m.Close()

	flvHeader := GetHeaderCompo()
	d.ReadHeader(flvHeader)
	m.WriteHeader(flvHeader)
	flvHeader.Put()

	for {
		select {
		case <-p.stop:
			return nil
		default:
			flvBody := new(TagCompo)

			err := d.ReadTag(flvBody)
			if err != nil {
				if err == io.EOF {
					return nil
				} else {
					return err
				}
			}

			if flvBody.ContainAVCSeqHeader() {
				if p.avcHeader == nil {
					p.avcHeader = make([]byte, len(flvBody.TagBodyRaw))
					copy(p.avcHeader, flvBody.TagBodyRaw)
				} else {
					if bytes.Compare(p.avcHeader, flvBody.TagBodyRaw) == 0 {
						flvBody.Free()
						continue
					} else {
						flvBody.Free()
						return fmt.Errorf("video(name = %s) AVC sequence header changed", out)
					}
				}
			}
			m.WriteTag(flvBody)
			flvBody.Free()
		}
	}

}
