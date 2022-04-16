package flv

import (
	"encoding/binary"
)

type HeaderCompo struct {
	Raw [13]byte

	Header
}

type Header struct {
	Signature       [3]uint8
	Version         uint8
	TypeFlags       uint8
	DataOffset      uint32
	PreviousTagSize uint32
}

type TagCompo struct {
	TagHeaderRaw [11]byte
	TagBodyRaw   []byte

	TagHeader
	// // equal to TagBodyRaw [0 : len-4]
	// TagData         []byte
	PreviousTagSize uint32
}

type TagHeader struct {
	TagType           uint8
	DataSize          [3]uint8
	Timestamp         [3]uint8
	TimestampExtended uint8
	StreamID          [3]uint8
}

func (this *TagHeader) GetDataSize() uint32 {
	dataSize := binary.BigEndian.Uint32(append([]uint8{0}, []uint8(this.DataSize[:])...))
	// log.Println("data size: ", dataSize)
	return dataSize
}

func (this *TagHeader) GetTimestamp() uint32 {
	realTimestamp := binary.BigEndian.Uint32(append([]uint8{this.TimestampExtended}, []uint8(this.Timestamp[:])...))
	// log.Printf("start time: %s", time.Unix(int64(fixedTimstamp/1000), 0))
	return realTimestamp
}
