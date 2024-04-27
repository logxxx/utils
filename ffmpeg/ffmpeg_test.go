package ffmpeg

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGenePreviewVideo(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	filePath := "H:\\ffmpeg_test\\转换失败\\Indie 轻奢女鞋_哇太赞了👍🏻这样的鞋还有谁不爱呀？.mp4"
	err := GenePreviewVideo(filePath, filePath+".thumb.mp4")
	if err != nil {
		panic(err)
	}

}

func TestGenePreviewVideo2(t *testing.T) {
	filePath := "H:\\ffmpeg_test\\特殊字符\\Indie 轻奢女鞋_哇太赞了👍🏻这样的鞋还有谁不爱呀？.mp4"
	newPath := strings.ReplaceAll(filePath, "👍🏻", "")
	err := os.Rename(filePath, newPath)
	if err != nil {
		panic(err)
	}
}

func TestGenePreviewVideo3(t *testing.T) {
	filePath := "H:\\ffmpeg_test\\特殊字符\\Indie 轻奢女鞋_哇太赞了👍🏻这样的鞋还有谁不爱呀？.mp4"
	//filePath = "H:\\ffmpeg_test\\特殊字符\\1.mp4"
	output, err := exec.Command("attrib", filePath).CombinedOutput()
	t.Logf("output:%v err:%v", string(output), err)
}
