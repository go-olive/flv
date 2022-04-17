package flv

import (
	"encoding/binary"
	"math"
)

var (
	sig = [3]uint8{0x46, 0x4c, 0x56}
	ver = uint8(0x01)
	do  = uint32(9)

	frameRate           = 15.0
	videoFrameInterval  = math.Ceil(1000 / frameRate)
	soundSampleInterval = math.Ceil(1000 / 44)
)

type HeaderCompo struct {
	Raw [13]byte

	Header
}

type Header struct {
	Signature [3]uint8
	Version   uint8
	// the below four vars combine as a var TypeFlags
	// TypeFlagsReserved1 [5]uint1
	// TypeFlagsAudio     uint1
	// TypeFlagsReserved2 uint1
	// TypeFlagsVideo     uint1
	TypeFlags       uint8
	DataOffset      uint32
	PreviousTagSize uint32
}

func (this *Header) Valid() bool {
	return this.Signature == sig &&
		this.Version == ver &&
		this.DataOffset == do
}

func (this *Header) HasAudio() bool {
	return this.TypeFlags&uint8(0b00000100) != 0
}
func (this *Header) HasVedio() bool {
	return this.TypeFlags&uint8(0b00000001) != 0
}

type TagCompo struct {
	TagHeaderRaw [11]byte
	TagBodyRaw   []byte

	TagHeader
	// // equal to TagBodyRaw [0 : len-4]
	// TagData         []byte
	PreviousTagSize uint32
}

func (this *TagCompo) Free() {
	PutBytes(this.TagBodyRaw)
}

type TagHeader struct {
	TagType byte
	// 3 bytes
	DataSize []byte
	// 4 bytes
	Timestamp         []byte
	TimestampExtended byte
	// 3 bytes
	StreamID []byte
}

func (this *TagHeader) GetDataSize() uint32 {
	dataSize := binary.BigEndian.Uint32(append([]uint8{0}, []uint8(this.DataSize[:])...))
	// log.Println("data size: ", dataSize)
	return dataSize
}

func (this *TagHeader) GetTimestamp() uint32 {
	if len(this.Timestamp) != 3 {
		return 0
	}
	realTimestamp := binary.BigEndian.Uint32(append([]uint8{this.TimestampExtended}, []uint8(this.Timestamp[:])...))
	// log.Printf("start time: %s", time.Unix(int64(fixedTimstamp/1000), 0))
	return realTimestamp
}

func (this *TagCompo) ContainAVCSeqHeader() bool {
	if this.TagHeader.TagType != 9 {
		return false
	}
	if len(this.TagBodyRaw) < 2 {
		return false
	}
	if (this.TagBodyRaw[0] & 15) != 7 {
		return false
	}
	if this.TagBodyRaw[1] == 0 {
		return true

	}
	return false
}
