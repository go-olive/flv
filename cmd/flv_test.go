package main

import (
	"io"
	"os"
	"testing"

	"github.com/go-olive/flv"
)

func BenchmarkFlvFileNone(b *testing.B) {
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		fix, _ := os.Create("fix_none.flv")
		m, _ := flv.NewMuxer(fix)
		defer m.Close()
		f, _ := os.Open("good.flv")
		d, _ := flv.NewDemuxer(f)
		defer d.Close()

		b.StartTimer()
		flvHeader := new(flv.HeaderCompo)
		d.ReadHeader(flvHeader)
		m.WriteHeader(flvHeader)
		for {
			flvBody := new(flv.TagCompo)
			err := d.ReadTag(flvBody)
			if err != nil {
				if err == io.EOF {
					println("success")
					break
				} else {
					println(err.Error())
					break
				}
			}
			m.WriteTag(flvBody)
		}
		b.StopTimer()
	}
}

// func BenchmarkHandlePool(b *testing.B) {
// 	b.StopTimer()

// 	for i := 0; i < b.N; i++ {

// 		fix, _ := os.Create("fix2.flv")
// 		m, _ := flv.NewMuxer(fix)

// 		f, _ := os.Open("good.flv")
// 		d, _ := flv.NewDemuxer(f)

// 		b.StartTimer()

// 		flvHeader := flv.GetHeaderCompo()
// 		d.ReadHeader(flvHeader)
// 		m.WriteHeader(flvHeader)
// 		flvHeader.Put()
// 		for {
// 			flvBody := flv.GetTagCompo()

// 			err := d.ReadTag(flvBody)
// 			if err != nil {
// 				if err == io.EOF {
// 					println("read to finish")
// 					break
// 				} else {
// 					println(err.Error())
// 					break
// 				}
// 			}
// 			m.WriteTag(flvBody)
// 			flvBody.Put()
// 		}
// 		b.StopTimer()
// 	}
// }

func BenchmarkFlvFileSlab(b *testing.B) {
	// var cnt uint8

	b.StopTimer()

	for i := 0; i < b.N; i++ {
		// cnt = 0

		b.StartTimer()
		fix, _ := os.Create("fix_slab.flv")
		m, _ := flv.NewMuxer(fix)
		defer m.Close()
		f, _ := os.Open("good.flv")
		d, _ := flv.NewDemuxer(f)
		defer d.Close()

		b.StartTimer()
		flvHeader := new(flv.HeaderCompo)
		d.ReadHeader(flvHeader)
		m.WriteHeader(flvHeader)
		for {
			flvBody := new(flv.TagCompo)

			err := d.ReadTag(flvBody)
			if err != nil {
				if err == io.EOF {
					println("success")
					break
				} else {
					println(err.Error())
					break
				}
			}
			m.WriteTag(flvBody)
			// if flvBody.ContainAVCSeqHeader() {
			// 	cnt++
			// 	println("avc header cnt: ", cnt)
			// }
			flvBody.Free()
		}
		b.StopTimer()
	}
}
