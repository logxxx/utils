package media_test

import (
	"github.com/logxxx/utils/media"
	"testing"
)

func TestTryCut(t *testing.T) {
	//1分13 ->2分15
	path := "N:\\source\\日常\\朋友圈\\25.mp4"
	media.CutVideo(path, 0, 6)
}
