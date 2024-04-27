package ffmpeg

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGenePreviewVideo(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	filePath := "H:\\ffmpeg_test\\转换失败\\1.mp4"
	err := GenePreviewVideo(filePath, filePath+".thumb.mp4")
	if err != nil {
		panic(err)
	}
}
