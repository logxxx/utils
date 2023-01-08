package media

import (
	"github.com/logxxx/utils/log"
	"image/jpeg"
	"os"
)

func CompressPic(imagePath string) (result string) {

	result = imagePath

	//需要压缩
	imgFile, err := os.Open(imagePath)
	if err != nil {
		log.Errorf("CompressPic os.Open err:%v req:%v", err, imagePath)
		return
	}
	defer imgFile.Close()

	fileInfo, _ := imgFile.Stat()
	if fileInfo.Size() < 1024*1024/2 {
		log.Infof("CompressPic 图片<500k，不需要压缩:%vkb", fileInfo.Size()/1024/2)
		return
	}
	log.Infof("CompressPic 图片>500k，需要压缩:%vkb", fileInfo.Size()/1024/2)

	jpgimg, err := jpeg.Decode(imgFile)
	if err != nil {
		log.Errorf("CompressPic Decode err:%v", err)
		return
	}

	//保存到新文件中
	newPath := imagePath + ".thumb"
	newfile, err := os.Create(newPath)
	if err != nil {
		log.Errorf("CompressPic Create err:%v", err)
		return
	}
	defer newfile.Close()

	// &jpeg.Options{Quality: 10} 图片压缩质量
	err = jpeg.Encode(newfile, jpgimg, &jpeg.Options{Quality: 30})
	if err != nil {
		log.Errorf("CompressPic Encode err:%v", err)
		return
	}

	newFileInfo, _ := newfile.Stat()
	log.Infof("CompressPic result: %vkb => %vkb", fileInfo.Size()/1024, newFileInfo.Size()/1024)

	return newPath
}
