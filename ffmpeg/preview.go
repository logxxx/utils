package ffmpeg

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/runutil"
	log "github.com/sirupsen/logrus"

	"os"
	"path/filepath"
	"strings"
)

type GenePreviewVideoSliceOpt struct {
	ToPath      string
	SegNum      int
	SegDuration int
	SkipStart   int
	SkipEnd     int
}

func GenePreviewVideoSlice(filePath string, fn func(vInfo *VideoFile) GenePreviewVideoSliceOpt) (resp string, err error) {

	logger := log.WithField("func_name", "GenePreviewVideoSlice").WithField("filePath", filePath)

	ffprobe := FFProbe("ffprobe")
	video, err := ffprobe.NewVideoFile(filePath)
	if err != nil {
		logger.Errorf("GenePreviewVideoSlice NewVideoFile err:%v", err)
		return
	}

	logger.Debugf("all duration:%v", video.Duration)

	opt := fn(video)

	cutPoints := getCutPoints(int(video.Duration), opt.SegNum, opt.SegDuration, opt.SkipStart, opt.SkipEnd)
	logger.Debugf("cutPoints:%v", cutPoints)

	chunks := make([]string, 0)

	previewDir := filepath.Join(filepath.Dir(filePath), fmt.Sprintf("_ffmpegpreview_%v", utils.MD5(filePath)))

	defer func() {
		runutil.RunSafe(func() error {

			logger.Debugf("remove preview dir:%v", previewDir)
			err := os.RemoveAll(previewDir)
			if err != nil {
				logger.Debugf("os.Remove err:%v path:%v", err, previewDir)
			} else {
				logger.Debugf("remove preview dir succ")
			}
			return nil
		})

	}()

	w, h := getPreviewWH(video)

	for i, point := range cutPoints {
		logger.Debugf("genePreviewVideoChunk %v/%v %v~%v", i+1, len(cutPoints), point, point+opt.SegDuration)
		chunk, err := genePreviewVideoChunk(filePath, previewDir, point, point+opt.SegDuration, w, h)
		if err != nil {
			logger.Errorf("GenePreviewVideo genePreviewVideoChunk err:%v", err)
			return "", err
		}
		chunks = append(chunks, chunk)
	}
	logger.Debugf("chunks:%v", chunks)

	mergedPath, err := mergeChunks(filePath, previewDir, chunks, opt.ToPath)
	if err != nil {
		logger.Errorf("GenePreviewVideo mergeChunks err:%v", err)
		return "", err
	}
	return mergedPath, nil

}

func getPreviewWH(v *VideoFile) (w int, h int) {
	min := 480
	w = v.Width
	h = v.Height
	if w > h { //宽
		for {
			if w <= min {
				return
			}
			w /= 2
			h /= 2
		}
	}
	if w < h { //长视频
		for {
			if h <= min {
				return
			}
			w /= 2
			h /= 2
		}
	}
	return
}

func mergeChunks(sourcePath string, previewDir string, chunks []string, toPath string) (string, error) {
	contactFile := filepath.Join(previewDir, fmt.Sprintf("ffmpeg_concat.txt"))
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

	for segmentDuration*segmentNum > videoDuration {
		if segmentDuration > 3 {
			segmentDuration--
		} else {
			break
		}
	}

	for segmentDuration*segmentNum > videoDuration {
		if segmentNum > 1 {
			segmentNum--
		} else {
			break
		}
	}

	for skipStart+skipEnd+segmentNum*segmentDuration > videoDuration {
		if skipEnd > 0 {
			skipEnd--
		} else {
			break
		}
	}

	for skipStart+skipEnd+segmentNum*segmentDuration > videoDuration {
		if skipStart > 0 {
			skipStart--
		} else {
			break
		}
	}

	log.Debugf("getCutPoints videoDuration:%v segNum:%v segDur:%v skipStart:%v skipEnd:%v", videoDuration, segmentNum, segmentDuration, skipStart, skipEnd)

	points := make([]int, 0)
	allPointNum := (videoDuration - skipStart - skipEnd) / (segmentDuration)
	step := allPointNum / segmentNum
	log.Debugf("getCutPoints allPointNum:%v step:%v", allPointNum, step)
	for i := 1; i <= segmentNum; i++ {
		points = append(points, skipStart+i*step*segmentDuration)
	}
	return points
}

func genePreviewVideoChunk(sourcePath, previewDir string, fromSec, toSec int, w, h int) (string, error) {

	command := "ffmpeg -y -ss %v -to %v -i %v -vf scale=%v:%v -pix_fmt yuv420p -profile:v high -level 4.2 -crf 21 -threads 4 -strict -2 %v"
	pureName, ext := getPureNameAndExt(sourcePath)
	outputFilePath := filepath.Join(previewDir, fmt.Sprintf("ffmpegtrunk_%v_%v~%vs%v", pureName, fromSec, toSec, ext))
	os.MkdirAll(filepath.Dir(outputFilePath), 0755)
	command = fmt.Sprintf(command, fromSec, toSec, sourcePath, w, h, outputFilePath)
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
