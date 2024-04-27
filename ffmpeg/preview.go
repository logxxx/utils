package ffmpeg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func GenePreviewVideo(filePath string, toPath string) error {

	scale := "480:-1"
	fpb := FFProbe("ffprobe")
	vInfo, err := fpb.NewVideoFile(filePath)
	if err == nil {
		log.Infof("height:%v width:%v", vInfo.Height, vInfo.Width)
		if vInfo.Height > vInfo.Width { //竖屏
			scale = "-1:480"
		}
	} else {
		log.Errorf("GenePreviewVideo NewVideoFile err:%v", err)
	}

	command := "ffmpeg -y -i %v -to 15 -vf scale=%v -pix_fmt yuv420p -level 4.2 -crf 30 -threads 8 -strict -2 %v"
	command = fmt.Sprintf(command, filePath, scale, toPath)
	output, err := runCommand(command)
	log.Debugf("GenePreviewVideo command:%v output:%v err:%v", command, string(output), err)
	if err != nil {
		return err
	}
	return nil
}
