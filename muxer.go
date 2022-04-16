package flv

import (
	"encoding/binary"
	"io"
)

type Muxer interface {
	WriteHeader(*HeaderCompo) (err error)
	WriteTag(*TagCompo) (err error)
	Close() error
}

func NewMuxer(w io.WriteCloser) (Muxer, error) {
	return &muxer{
		w: w,
	}, nil
}

type muxer struct {
	w io.WriteCloser

	offset uint32
}

func (m *muxer) WriteHeader(c *HeaderCompo) (err error) {
	_, err = m.w.Write(c.Raw[:])
	return
}

func (m *muxer) WriteTag(c *TagCompo) (err error) {
	if m.offset == 0 {
		m.offset = c.GetTimestamp()
	}
	fixedTimstamp := c.GetTimestamp() - m.offset
	modified := make([]byte, 4)
	binary.BigEndian.PutUint32(modified, fixedTimstamp)

	c.TagHeaderRaw[7] = modified[0]
	copy(c.TagHeaderRaw[4:7], modified[1:4])

	_, err = m.w.Write(c.TagHeaderRaw[:])
	if err != nil {
		return
	}
	_, err = m.w.Write(c.TagBodyRaw)
	return
}

func (m *muxer) Close() error {
	return m.w.Close()
}
