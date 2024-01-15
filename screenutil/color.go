package screen

import (
	"fmt"
	"github.com/logxxx/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

type CheckColorFunc func(r, g, b uint8) bool

var (
	IsYellow = func(r, g, b uint8) bool {
		//深黄:247,162,24
		//亮黄:255,255,130
		if r >= 247 && g >= 162 && b <= 130 {
			return true
		}
		return false
	}

	IsGreen = func(r, g, b uint8) bool {
		//深绿:24, 178, 148
		//亮绿:104,249,229
		if r <= 104 && g >= 178 && b >= 148 {
			return true
		}
		return false
	}
)

func CheckColor(x1, y1, x2, y2 int, isDebug bool, fns ...CheckColorFunc) bool {

	if len(fns) == 0 {
		panic("CheckColor err: len(fns) == 0")
	}

	img, err := ShotRectOrig(x1, y1, x2, y2)
	if err != nil {
		return false
	}

	if isDebug {
		filename := fmt.Sprintf("./asset/%v.jpg", utils.FormatTimeSafe(time.Now()))
		err = SaveToLocal(img, filename)
		if err != nil {
			return false
		}
	}

	totalCount := 0

	//接取任务(绿)和领奖(黄) 都在同一个位置。

	results := make([]int, len(fns))

	for i := 0; i < x2-x1; i++ {
		for j := 0; j < y2-y1; j++ {
			totalCount++
			rawR, rawG, rawB, _ := img.At(i, j).RGBA()
			r := uint8(rawR) //转成255
			g := uint8(rawG) //转成255
			b := uint8(rawB) //转成255
			//log.Infof("[%v, %v] (%v,%v,%v,%v)", i+x1, j+y1, r, g, b)

			for i, fn := range fns {
				ok := fn(r, g, b)
				if ok {
					results[i]++
					continue
				}
			}
		}
	}

	for i, result := range results {
		prob := float64(result*1000/totalCount/100) / 10
		if isDebug {
			log.Infof("idx:%v result:%v/%v prob:%v", i, result, totalCount, prob)
		}
		if prob >= 0.5 {
			return true
		}
	}

	return false
}
