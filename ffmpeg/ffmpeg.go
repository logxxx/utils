package ffmpeg

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GeneScreenShot(sourcePath string, point int) (string, error) {
	pureName, _ := fileutil.GetPureNameAndExt(sourcePath)
	outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("%v_第%v秒.jpg", pureName, point))
	//ffmpeg.exe -ss 10 -i possible.mkv -y -f image2 -t 0.01 0.jpg
	command := fmt.Sprintf("ffmpeg -ss %v -i %v -y -f image2 -t 0.01 %v", point, sourcePath, outputPath)
	err := runCommand(command)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

func runCommand(command string) error {
	log.Infof("runCommand:%v", command)
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
