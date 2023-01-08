package media

import (
	"golang.org/x/image/webp"
	"image/jpeg"
	"os"
	"strings"
)

func Webp2Jpg(path string) (string, error) {

	webpFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer webpFile.Close()

	m, err := webp.Decode(webpFile)
	if err != nil {
		return "", err
	}
	destFileName := strings.ReplaceAll(path, ".webp", ".jpg")
	destFile, err := os.OpenFile(destFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer destFile.Close()
	err = jpeg.Encode(destFile, m, &jpeg.Options{Quality: 100})
	if err != nil {
		return "", err
	}
	return destFileName, nil
}
