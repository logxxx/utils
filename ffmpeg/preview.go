package ffmpeg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func GenePreviewVideo(filePath string, toPath string) error {
	command := "ffmpeg -y -i %v -to 15 -vf scale=640:-2 -pix_fmt yuv420p -level 4.2 -crf 21 -threads 4 -strict -2 %v"
	command = fmt.Sprintf(command, filePath, toPath)
	output, err := runCommand(command)
	log.Debugf("GenePreviewVideo command:%v output:%v err:%v", command, string(output), err)
	if err != nil {
		return err
	}
	return nil
}
