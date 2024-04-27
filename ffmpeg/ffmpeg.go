package ffmpeg

import (
	"fmt"
	"github.com/logxxx/utils/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GeneScreenShot(sourcePath string, point int) (string, error) {
	pureName, _ := getPureNameAndExt(sourcePath)
	outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("%v_第%v秒.jpg", pureName, point))
	//ffmpeg.exe -ss 10 -i possible.mkv -y -f image2 -t 0.01 0.jpg
	command := fmt.Sprintf("ffmpeg -ss %v -i %v -y -f image2 -t 0.01 %v", point, sourcePath, outputPath)
	_, err := runCommand(command)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

func getPureNameAndExt(sourcePath string) (string, string) {
	baseName := filepath.Base(sourcePath)
	ext := filepath.Ext(baseName)
	pureName := strings.TrimSuffix(baseName, ext)
	return pureName, ext
}

func runCommand(command string) (output []byte, err error) {
	log.Infof("runCommand:%v", command)
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	return
}
