package flv

import "sync"

var (
	headerCompoPool = sync.Pool{
		New: func() any { return new(HeaderCompo) },
	}
	tagCompoPool = sync.Pool{
		New: func() any {
			return new(TagCompo)
		},
	}
)

func GetHeaderCompo() *HeaderCompo {
	return headerCompoPool.Get().(*HeaderCompo)
}

func (this *HeaderCompo) Put() {
	headerCompoPool.Put(this)
}

func GetTagCompo() *TagCompo {
	return tagCompoPool.Get().(*TagCompo)
}

func (this *TagCompo) Put() {
	if cap(this.TagBodyRaw) > 64<<10 {
		return
	}
	this.TagBodyRaw = this.TagBodyRaw[:0]
	tagCompoPool.Put(this)
}
