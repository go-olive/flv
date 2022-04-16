package flv

import (
	"sync"

	"github.com/funny/slab"
)

var (
	headerCompoPool = sync.Pool{
		New: func() any { return new(HeaderCompo) },
	}
	bytesPool = slab.NewChanPool(
		16,        // The smallest chunk size is 16B.
		256*1024,  // The largest chunk size is 256KB.
		2,         // Power of 2 growth in chunk size.
		1024*1024, // Each slab will be 1MB in size.
	)
)

func GetHeaderCompo() *HeaderCompo {
	return headerCompoPool.Get().(*HeaderCompo)
}

func (this *HeaderCompo) Put() {
	headerCompoPool.Put(this)
}

func PutBytes(slice []byte) {
	bytesPool.Free(slice)
}

func GetBytes(s int) []byte {
	return bytesPool.Alloc(s)
}
