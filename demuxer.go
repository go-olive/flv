package flv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Demuxer interface {
	ReadHeader(*HeaderCompo) error
	ReadTag(*TagCompo) error
	Close() error
}

func NewDemuxer(r io.ReadCloser) (Demuxer, error) {
	return &demuxer{
		r: r,
	}, nil
}

type demuxer struct {
	r io.ReadCloser
}

func (d *demuxer) ReadHeader(c *HeaderCompo) (err error) {
	if _, err = io.ReadFull(d.r, c.Raw[:]); err != nil {
		return
	}

	reader := bytes.NewReader(c.Raw[:])
	if err = binary.Read(reader, binary.BigEndian, &c.Header); err != nil {
		return
	}
	if err = binary.Read(reader, binary.BigEndian, &c.PreviousTagSize); err != nil {
		return
	}

	return nil
}

func (d *demuxer) ReadTag(c *TagCompo) (err error) {
	if _, err = io.ReadFull(d.r, c.TagHeaderRaw[:]); err != nil {
		return
	}

	c.TagHeader.TagType = c.TagHeaderRaw[0]
	c.TagHeader.DataSize = c.TagHeaderRaw[1:4]
	c.TagHeader.Timestamp = c.TagHeaderRaw[4:7]
	c.TagHeader.TimestampExtended = c.TagHeaderRaw[7]
	c.TagHeader.StreamID = c.TagHeaderRaw[8:11]

	if len(c.TagBodyRaw) < int(c.GetDataSize()+4) {
		c.TagBodyRaw = make([]byte, c.GetDataSize()+4)
	}

	if _, err = io.ReadAtLeast(io.LimitReader(d.r, int64(c.GetDataSize()+4)), c.TagBodyRaw, int(c.GetDataSize()+4)); err != nil {
		return
	}

	c.PreviousTagSize = binary.BigEndian.Uint32(c.TagBodyRaw[c.GetDataSize() : c.GetDataSize()+4])
	if c.GetDataSize()+11 != c.PreviousTagSize {
		return fmt.Errorf("Data size incorrect, got: %d, want: %d.", c.PreviousTagSize, c.GetDataSize()+11)
	}

	return nil

}

func (d *demuxer) Close() error {
	return d.r.Close()
}
