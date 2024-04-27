package ffmpeg

import "testing"

func TestGenePreviewVideo(t *testing.T) {
	filePath := "H:\\output_xhs\\20240316\\测试视频\\Gandy___这个背部错误，浪费你很多时间.mp4"
	err := GenePreviewVideo(filePath, filePath+".thumb.mp4")
	if err != nil {
		panic(err)
	}
}
