package main

import (
	"github.com/go-olive/flv"
)

func main() {
	// t, err := tv.Snap(tv.NewRoomUrl("https://www.bilibili.com/512709"), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// t.Refresh()
	// fmt.Println(t)
	// s, _ := t.StreamUrl()
	s := ""
	p := flv.NewParser()
	if err := p.Parse(s, "hy.flv"); err != nil {
		println(err.Error())
	}
}
