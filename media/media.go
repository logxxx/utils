package media

import (
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"path/filepath"
	"strings"
)

var (
	AllowVideoExts = []string{".wmv", ".asf", ".asx", ".rm", ".rmvb", ".mpg", ".mpeg", ".mpe", ".3gp", ".mov", ".mp4", ".m4v", ".avi", ".mkv", ".flv", ".f4v", ".vob", ".ts", ".wm", ".m1v", ".m2v", ".mpv", ".mpv2", ".mp2v", ".tp", ".tpr", ".ifo", ".ogm", ".ogv", ".m4p", ".v4b", ".3gpp", ".3g2", ".3gp2", ".ram", ".rpm", ".qt", ".nsv", ".dpg", ".m2ts", ".m2t", ".mts", ".k3g", ".skm", ".evo", ".nsr", ".amv", ".divx", ".webm", ".wtv", ".f4v", ".mfx", ".h264", ".trt", ".m2p"}
	AllowImageExts = []string{"bmp", "jpg", "jpeg", "webp", "png", "tif", "gif", "pcx", "tga", "exif", "fpx", "svg", "psd", "cdr", "pcd", "dxf", "ufo", "eps", "ai", "raw", "WMF", "webp", "avif", "apng"}
)

func IsVideo(fileName string) bool {
	if fileutil.IsDir(fileName) {
		return false
	}
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)

	return utils.Contains(ext, AllowVideoExts)

}

func IsImage(fileName string) bool {
	if fileutil.IsDir(fileName) {
		return false
	}
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)

	return utils.Contains(ext, AllowImageExts)
}
