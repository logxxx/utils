package main

import (
	"fmt"
	"github.com/logxxx/utils/media"
	log "github.com/sirupsen/logrus"
	"time"
)

func main1() {

	path := `D:\mytest\mywork\myweibo\backend\source\蛋壳-安利协会(1876856920)\20220716_1.mp4`

	newPath := media.TryCut(path)
	log.Printf("newPath:%v", newPath)

}

func main2() {
	path := `C:\Users\hehe\Desktop\test.mp4`
	newPath, err := media.Reformat(path)
	if err != nil {
		panic(err)
	}
	log.Printf("newPath:%v", newPath)
}

func main() {

	for i := 0; i <= 8; i++ {
		st := time.Now()
		path := fmt.Sprintf(`D:\mytest\ae\1\source\%v.webp`, i)
		resp, err := media.Webp2Jpg(path)
		if err != nil {
			panic(err)
		}
		log.Printf("resp:%v useTime:%v", resp, time.Now().Sub(st).String())
	}
}
