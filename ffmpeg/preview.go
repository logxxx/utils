package ffmpeg

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"os"
	"path/filepath"
	"strings"
	"time"
)

type GenePreviewVideoSliceOpt struct {
	FilePath    string
	ToPath      string
	SegNum      int
	SegDuration int
	SkipStart   int
	SkipEnd     int
}

func GenePreviewVideoSlice(opt GenePreviewVideoSliceOpt) (resp string, err error) {

	ffprobe := FFProbe("ffprobe")
	video, err := ffprobe.NewVideoFile(opt.FilePath)
	if err != nil {
		log.Errorf("GenePreviewVideoSlice NewVideoFile err:%v", err)
		return
	}

	log.Debugf("all duration:%v", video.Duration)

	cutPoints := getCutPoints(int(video.Duration), opt.SegNum, opt.SegDuration, opt.SkipStart, opt.SkipEnd)
	log.Debugf("cutPoints:%v", cutPoints)

	height := video.Height
	width := video.Width

	min := 640
	if video.Height > video.Width { //竖屏

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

	chunks := make([]string, 0)

	defer func() {
		os.RemoveAll(filepath.Join(filepath.Dir(opt.FilePath), "_ffmpegpreview"))
		os.RemoveAll(filepath.Join(filepath.Dir(opt.FilePath), "_preview"))
	}()

	for i, point := range cutPoints {
		log.Debugf("genePreviewVideoChunk %v/%v %v~%v", i+1, len(cutPoints), point, point+opt.SegDuration)
		chunk, err := genePreviewVideoChunk(opt.FilePath, point, point+opt.SegDuration, scale)
		if err != nil {
			log.Errorf("GenePreviewVideo genePreviewVideoChunk err:%v", err)
			return "", err
		}
		chunks = append(chunks, chunk)
	}
	log.Printf("chunks:%v", chunks)

	mergedPath, err := mergeChunks(opt.FilePath, chunks, opt.ToPath)
	if err != nil {
		log.Printf("GenePreviewVideo mergeChunks err:%v", err)
		return "", err
	}
	return mergedPath, nil

}

func mergeChunks(sourcePath string, chunks []string, toPath string) (string, error) {
	contactFile := filepath.Join(filepath.Dir(sourcePath), "_preview", fmt.Sprintf("ffmpeg_concat_%v.txt", time.Now().UnixNano()))
	os.MkdirAll(filepath.Dir(contactFile), 0755)
	content := ""
	for _, chunk := range chunks {
		content += fmt.Sprintf("file '%v'\n", chunk)
	}
	err := os.WriteFile(contactFile, []byte(content), 0755)
	if err != nil {
		log.Errorf("mergeChunks WriteFile err:%v content:%v", err, content)
		return "", err
	}
	pureName, _ := getPureNameAndExt(sourcePath)
	if toPath == "" {
		toPath = filepath.Join(toPath, fmt.Sprintf("%v_preview.mp4", pureName))
	}
	os.MkdirAll(filepath.Dir(toPath), 0755)

	command := fmt.Sprintf("ffmpeg -y -f concat -safe 0 -i %v %v", contactFile, toPath)
	_, err = runCommand(command)
	if err != nil {
		log.Errorf("mergeChunks runCommand err:%v command:%v", err, command)
		return "", err
	}
	return toPath, nil
}

func getCutPoints(videoDuration int, segmentNum int, segmentDuration int, skipStart, skipEnd int) []int {

	if segmentDuration <= 0 {
		segmentDuration = 5
	}

	if segmentNum <= 0 {
		segmentNum = 3
	}

	points := make([]int, 0)
	allPointNum := (videoDuration - skipStart - skipEnd) / (segmentDuration)
	step := allPointNum / segmentNum
	log.Debugf("getCutPoints allPointNum:%v step:%v", allPointNum, step)
	if allPointNum <= 0 {
		segmentNum = 3
		segmentDuration = 5
		skipStart = 0
		skipEnd = 0
		allPointNum = (videoDuration - skipStart - skipEnd) / (segmentDuration)
		step = allPointNum / segmentNum
	}
	for i := 1; i <= segmentNum; i++ {
		points = append(points, skipStart+i*step*segmentDuration)
	}
	return points
}

func genePreviewVideoChunk(sourcePath string, fromSec, toSec int, scale string) (string, error) {
	sourceDir := filepath.Dir(sourcePath)
	command := "ffmpeg -y -ss %v -to %v -i %v -vf scale=%v -pix_fmt yuv420p -profile:v high -level 4.2 -crf 21 -threads 4 -strict -2 %v"
	pureName, ext := getPureNameAndExt(sourcePath)
	outputFilePath := filepath.Join(sourceDir, "_ffmpegpreview", fmt.Sprintf("ffmpegtrunk_%v_%v~%vs%v", pureName, fromSec, toSec, ext))
	os.MkdirAll(filepath.Dir(outputFilePath), 0755)
	command = fmt.Sprintf(command, fromSec, toSec, sourcePath, scale, outputFilePath)
	_, err := runCommand(command)
	if err != nil {
		return "", err
	}
	return outputFilePath, nil
}

func getPureNameAndExt(sourcePath string) (string, string) {
	baseName := filepath.Base(sourcePath)
	ext := filepath.Ext(baseName)
	pureName := strings.TrimSuffix(baseName, ext)
	return pureName, ext
}

func GenePreviewVideo(filePath string, toPath string) error {

	fpb := FFProbe("ffprobe")
	vInfo, err := fpb.NewVideoFile(filePath)
	if err != nil {
		log.Errorf("GenePreviewVideo NewVideoFile err:%v", err)
		return err
	}
	//log.Infof("height:%v width:%v", vInfo.Height, vInfo.Width)

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

	command := `ffmpeg -y -i %v -to 15 -vf scale=%v -pix_fmt yuv420p -level 4.2 -crf 30 -threads 8 -strict -2 %v`
	command = fmt.Sprintf(command, filePath, scale, toPath)
	output, err := runCommand(command)
	log.Debugf("GenePreviewVideo command:%v output:%v err:%v", command, string(output), err)
	if err != nil {
		return err
	}
	return nil
}
