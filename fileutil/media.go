package fileutil

import (
	"github.com/logxxx/utils"
	"os"
	"path/filepath"
	"strings"
)

var (
	AllowVideoExts = []string{".wmv", ".asf", ".asx", ".rm", ".rmvb", ".mpg", ".mpeg", ".mpe", ".3gp", ".mov", ".mp4", ".m4v", ".avi", ".mkv", ".flv", ".f4v", ".vob", ".ts", ".wm", ".m1v", ".m2v", ".mpv", ".mpv2", ".mp2v", ".tp", ".tpr", ".ifo", ".ogm", ".ogv", ".m4p", ".v4b", ".3gpp", ".3g2", ".3gp2", ".ram", ".rpm", ".qt", ".nsv", ".dpg", ".m2ts", ".m2t", ".mts", ".k3g", ".skm", ".evo", ".nsr", ".amv", ".divx", ".webm", ".wtv", ".f4v", ".mfx", ".h264", ".trt", ".m2p"}
)

func IsVideo(fileName string) bool {

	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)

	if ext == ".ts" {
		return CheckTsFileIsVideo(fileName)
	}

	return utils.Contains(ext, AllowVideoExts)

}

func CheckTsFileIsVideo(filePath string) (is bool) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	header := make([]byte, 1)
	_, err = f.Read(header)
	if err != nil {
		return
	}
	return header[0] == 0x47
}
