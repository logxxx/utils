package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/logxxx/utils/log"
	"github.com/logxxx/utils/media"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetShowSize(input int64) string {
	size := float64(input)
	units := []string{"b", "kb", "MB", "GB", "TB"}
	unitIdx := 0
	for {
		if size < 1024 || unitIdx >= len(units) {
			break
		}
		size /= 1024
		unitIdx++
	}
	return fmt.Sprintf("%.2f%v", size, units[unitIdx])
}

func StartCorn(fn func() error, d time.Duration) {

	ticker := time.Tick(1 * time.Second)

	for {
		select {
		case <-ticker:
			ticker = time.Tick(d)

			if fn == nil {
				continue
			}

			err := fn()
			if err != nil {
				log.Errorf("startCorn fn err:%v", err)
			}

			log.Infof("cron next run time:%v d:%v", time.Now().Add(d).Format("2006/01/02 15:04:05"), d.String())

		}
	}

}

func AddSubfix(path string, subfix string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	pureFileName := strings.Trim(base, ext)
	pureFileName += subfix
	base = pureFileName + ext
	return filepath.Join(filepath.Dir(path), base)
}

func ContainsPath(input string, slice []string) bool {
	for _, elem := range slice {
		if filepath.Clean(input) == filepath.Clean(elem) {
			return true
		}
	}
	return false
}

func Contains(input string, slice []string) bool {
	for _, elem := range slice {
		if elem == input {
			return true
		}
	}
	return false
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func HasFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func GetFileSize(path string) (showSize string) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return "0kb"
	}
	size := fInfo.Size()
	format := "b"
	if size < 1024 {
		return fmt.Sprintf("%v%v", size, format)
	}
	format = "kb"
	size /= 1024
	if size < 1024 {
		return fmt.Sprintf("%v%v", size, format)
	}
	format = "mb"
	sizeFloat := float64(size) / 1024
	return fmt.Sprintf("%.2f%v", sizeFloat, format)
}

func JsonToString(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func ShortTitle(title string, limit ...int) string {
	return ShortName(TrimAt(TrimBracket(TrimUnderscore(EscapeEmoji(EscapeFileName(title))))), limit...)
}

func TrimAt(title string) string {
	startIdx := strings.Index(title, "@")
	if startIdx < 0 {
		return title
	}
	endIdx := strings.Index(title[startIdx:], " ")
	if endIdx < 0 {
		return title[:startIdx]
	}
	front := title[:startIdx]
	end := title[(startIdx + len(" ") + endIdx):]
	log.Printf("front:%v end:%v", front, end)
	return front + end
}

func TrimBracket(title string) string {
	startIdx := strings.Index(title, "【")
	if startIdx < 0 {
		return title
	}
	endIdx := strings.Index(title[startIdx:], "】")
	if endIdx < 0 {
		return title
	}
	front := title[:startIdx]
	end := title[(startIdx + len("【") + endIdx):]
	log.Printf("front:%v end:%v", front, end)
	return front + end
}

func TrimUnderscore(title string) string {
	idx := strings.Index(title, "_")
	if idx < 0 {
		return title
	}
	return title[idx+len("_"):]
}

func ShortName(input string, reqLimit ...int) string {

	runes := []rune(input)

	limit := 20
	if len(reqLimit) > 0 {
		limit = reqLimit[0]
	}

	if len(runes) <= limit {
		return input
	}

	return string(runes[:limit])
}

func EscapeFileName(input string) string {

	invalids := `~-.\/:*?"<>| `

	for _, invalid := range invalids {
		input = strings.ReplaceAll(input, string(invalid), "")
	}

	input = strings.ReplaceAll(input, "\n", "")

	input = strings.ReplaceAll(input, " ", "~")

	return input
}

func EscapeEmoji(input string) string {
	input = strings.ReplaceAll(input, `\ufffd`, "")
	for {
		startIdx := strings.Index(input, "[")
		if startIdx < 0 {
			return input
		}
		endIdx := strings.Index(input[startIdx:], "]")
		if endIdx < 0 {
			return input
		}
		input = input[:startIdx] + input[startIdx+endIdx+1:]

		if input == "" {
			input = fmt.Sprintf("%v", time.Now().Unix())
		}
	}
}

func SaveSourceVideo(path string) bool {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	destPah := filepath.Join(dir, "source")
	//if HasFile(destPah) {
	//	return false
	//}
	os.MkdirAll(destPah, 0666)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("SaveSourceVideo ReadFile err:%v path:%v", err, path)
		return false
	}
	err = ioutil.WriteFile(filepath.Join(dir, "source", base), content, 0666)
	if err != nil {
		log.Errorf("SaveSourceVideo WriteFile err:%v sourcePath:%v len(content):%v", err, filepath.Join(dir, "source", base), len(content))
		return false
	}

	return true
}

func IsSizeLargeThanMB(path string, delta int64) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	size := file.Size()
	log.Debugf("size:%.2fMB file:%v", float64(size)/1024/1024, path)
	return size > delta*1024*1024
}

func TrimVideo(downloadPath string, isTop, isBottom bool) error {

	videoInfo, err := media.GetMediaInfo(downloadPath)
	if err != nil {
		log.Errorf("GetMediaInfo err:%v", err)
		return err
	}

	tmpFile := downloadPath + ".mp4"

	top := 0
	bottom := 0
	if isTop {
		top = 100
	}
	if isBottom {
		bottom = 100
	}

	command := fmt.Sprintf("ffmpeg -i %v -vf crop=%v:%v:%v:%v %v -y", downloadPath, videoInfo.Width, videoInfo.Height-(top+bottom), 0, top, tmpFile)
	_, err = media.RunCmd(strings.Split(command, " "))
	if err != nil {
		log.Errorf("RunCmd err:%v", err)
		return err
	}

	os.Remove(downloadPath)
	os.Rename(tmpFile, downloadPath)

	return nil
}

func Extract(content string, begin string, end string) string {
	beginIdx := strings.Index(content, begin)
	if beginIdx < 0 {
		return ""
	}
	endIdx := strings.Index(content[beginIdx+len(begin):], end)
	if endIdx < 0 {
		return ""
	}
	return content[beginIdx+len(begin) : beginIdx+len(begin)+endIdx]
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetRandomOne(req []string) string {
	if len(req) <= 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return req[len(req)]
}
