package utils

import (
	"github.com/logxxx/utils/log"
	"testing"
)

func TestGetShowSize(t *testing.T) {

	req1 := int64(100 * 1024 * 1024)
	resp1 := GetShowSize(req1)
	t.Logf("resp1:%v", resp1)

}

func TestExtractAll(t *testing.T) {
	req := `        "live_photo": [
            "https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F003dwdargx086F50m58j0f0f0100abAw0k01.mov",
            "https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F002NuLrFgx086F50pczu0f0f0100c4040k01.mov"
        ],
`
	resp := ExtractAll(req, "https:", ".mov", true)
	log.Printf("resp:%v", resp)
}

func TestRemoveDuplicate(t *testing.T) {
	req := []string{
		"https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F003dwdargx086F50m58j0f0f0100abAw0k01.mov",
		"https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F002NuLrFgx086F50pczu0f0f0100c4040k01.mov",
		"https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F003dwdargx086F50m58j0f0f0100abAw0k01.mov",
		"https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F002NuLrFgx086F50pczu0f0f0100c4040k01.mov",
	}
	resp := RemoveDuplicate(req)
	log.Printf("resp(%v):%v", len(resp), resp)
}
