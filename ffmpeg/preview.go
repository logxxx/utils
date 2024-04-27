package ffmpeg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func GenePreviewVideo(filePath string, toPath string) error {

	fpb := FFProbe("ffprobe")
	vInfo, err := fpb.NewVideoFile(filePath)
	if err != nil {
		log.Errorf("GenePreviewVideo NewVideoFile err:%v", err)
		return err
	}
	log.Infof("height:%v width:%v", vInfo.Height, vInfo.Width)

	height := vInfo.Height
	width := vInfo.Width

	min := 640
	if vInfo.Height > vInfo.Width { //竖屏

		for {
			if height <= min {
				break
			}
			height /= 2
			width /= 2
		}

	} else {
		for {
			if width <= min {
				break
			}
			height /= 2
			width /= 2
		}
	}

	if width%2 != 0 {
		width -= 1
	}

	if height%2 != 0 {
		height -= 1
	}

	scale := fmt.Sprintf("%v:%v", width, height)

	command := `ffmpeg -y -i '%v' -to 15 -vf scale=%v -pix_fmt yuv420p -level 4.2 -crf 30 -threads 8 -strict -2 '%v'`
	command = fmt.Sprintf(command, filePath, scale, toPath)
	output, err := runCommand(command)
	log.Debugf("GenePreviewVideo command:%v output:%v err:%v", command, string(output), err)
	if err != nil {
		return err
	}
	return nil
}
