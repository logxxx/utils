package media

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logxxx/utils/log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type VideoInfo struct {
	DurationSec int
	Width       int
	Height      int
	Size        int64
}

type FFProbeResult struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	Format struct {
		FileName    string `json:"filename"`
		DurationStr string `json:"duration"`
		Size        string `json:"size"`
	}
}

func CutVideoFrontAndTail(path string, start, end int) error {

	videoInfo, err := GetMediaInfo(path)
	if err != nil {
		log.Errorf("CutVideoFrontAndTail GetMediaInfo err:%v", err)
		return err
	}

	if videoInfo.DurationSec-start-end < 0 {
		log.Printf("CutVideoFrontAndTail failed: dur:%v start:%v end:%v", videoInfo.DurationSec, start, end)
		return nil
	}

	newFile, err := CutVideo(path, start, videoInfo.DurationSec-end)
	if err != nil {
		log.Errorf("CutVideoFrontAndTail err:%v req1:%v req2:%+v", err, path, videoInfo)
		return err
	}
	os.Remove(path)
	os.Rename(newFile, path)

	return nil

}

func TryCut(path string) (result string) {

	result = path

	videoInfo, err := GetMediaInfo(path)
	if err != nil {
		return
	}

	if videoInfo.DurationSec >= 60 {

		start, end := getCutSec(videoInfo)

		output, err := CutVideo(path, start, end)
		if err != nil {
			log.Errorf("TryCut CutVideo err:%v req1:%v req2:%+v", err, path, videoInfo)
			return
		}
		result = output
	}

	return

}

func TrimVideo(path string) {

}

func GetMediaInfo(path string) (videoInfo *VideoInfo, err error) {

	args := []string{"ffprobe",
		"-show_format",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height", "-of", "json",
		path}
	result, err := RunCmd(args)
	if err != nil {
		log.Errorf("GetMediaInfo RunCmd err:%v", err)

	}

	//log.Infof("videoInfo:%v", result)

	resultObj := &FFProbeResult{}

	err = json.Unmarshal([]byte(result), resultObj)
	if err != nil {
		log.Errorf("GetMediaInfo RunCmd err:%v", err)

	}

	dFloat, err := strconv.ParseFloat(resultObj.Format.DurationStr, 32)
	if err != nil {
		log.Errorf("GetMediaInfo ParseFloat err:%v req:%v", err, resultObj.Format.DurationStr)

	}

	size, err := strconv.ParseInt(resultObj.Format.Size, 10, 64)
	if err != nil {
		log.Errorf("GetMediaInfo ParseInt size err:%v req:%v", err, resultObj.Format.Size)

	}

	width := 0
	height := 0

	if len(resultObj.Streams) > 0 {
		width = resultObj.Streams[0].Width
		height = resultObj.Streams[0].Height
	}

	videoInfo = &VideoInfo{
		DurationSec: int(dFloat),
		Width:       width,
		Height:      height,
		Size:        size,
	}

	return

}

func getCutOutputPath(path string) string {
	baseName := filepath.Base(path)
	outputName := baseName + ".thumb"
	output := filepath.Join(filepath.Dir(path), outputName)
	return output
}

func getResizedArgs(videoInfo *VideoInfo) []string {
	if videoInfo.Width <= 0 || videoInfo.Height <= 0 {
		return nil
	}

	return []string{"-vf", fmt.Sprintf("scale=%v:%v", videoInfo.Width/3*2, videoInfo.Height/3*2)}
}

func Reformat(path string) (string, error) {
	dir := filepath.Dir(path)
	baseName := filepath.Base(path)
	newFile := filepath.Join(dir, baseName+".mp4")
	rawCmd := fmt.Sprintf("ffmpeg -i %v -vcodec h264 -preset fast -b:v 2000k %v", path, newFile)
	args := strings.Split(rawCmd, " ")

	result, err := RunCmd(args)
	if err != nil {
		log.Errorf("CutVideo RunCmd err:%v", err)
		return "", err
	}
	_ = result
	//log.Infof("result:%v", result)

	return newFile, nil

}

func CutVideo(path string, start, end int) (string, error) {

	output := getCutOutputPath(path)

	formatStart := parseSecToStr(start)
	formatEnd := parseSecToStr(end)

	args := []string{"ffmpeg", "-ss", formatStart, "-to", formatEnd, "-i", path, "-y",
		"-f", "mp4",
		"-vcodec", "copy", "-acodec", "copy",
		"-q:v", "1",
		//"-c:v", "libx265", "-x265-params", "crf=18", //说是无损压缩，加上看不出来啥区别，视频尺寸还更大了...
	}
	//args = append(args, getResizedArgs(videoInfo)...) resize并不能减少太多体积
	args = append(args, output)

	log.Infof("Cut args:%v", args)

	result, err := RunCmd(args)
	if err != nil {
		log.Errorf("CutVideo RunCmd err:%v", err)
		return "", err
	}
	_ = result
	//log.Infof("result:%v", result)

	return output, nil

}

func RunCmd(input []string) (string, error) {
	if len(input) <= 0 {
		return "", errors.New("empty input")
	}

	log.Printf("RunCmd input:%v", input)

	command := input[0]
	args := input[1:]
	//switch command {
	//case "ffmpeg", "ffprobe", "convert":
	//default:
	//	return "", errors.New("unknown command")
	//}

	cmd := exec.Command(command, args...)

	out := bytes.NewBuffer(nil)
	cmd.Stdout = out

	err := cmd.Run()
	if err != nil {
		log.Printf("RunCmd err:%v output:%v", err, out.String())
		return "", err
	}

	return out.String(), nil

}

func parseSecToStr(totalSec int) string {

	min := totalSec / 60
	sec := totalSec % 60

	return fmt.Sprintf("00:%02d:%02d", min, sec)

}

func getCutSec(videoInfo *VideoInfo) (int, int) {
	//取中间1min
	midSec := videoInfo.DurationSec / 2
	startSec := midSec - 30
	if startSec < 0 {
		startSec = 0
	}
	//startSec -= 5
	//if startSec < 0 {
	//	startSec = 0
	//}

	endSec := midSec + 8
	if endSec > videoInfo.DurationSec {
		endSec = videoInfo.DurationSec
	}
	//endSec += 5
	//if endSec > totalSec {
	//	endSec = totalSec
	//}

	return startSec, endSec
}
