package ffmpeg

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGenePreviewVideo(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	filePath := "H:\\output_xhs\\20240316\\测试竖屏\\1.mp4"
	err := GenePreviewVideo(filePath, filePath+".111.thumb.mp4")
	if err != nil {
		panic(err)
	}
}
