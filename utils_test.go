package utils

import (
	log "github.com/sirupsen/logrus"
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

func TestExtract(t *testing.T) {
	req := `您的ID:wxid_ykfchhql9nw522

目前账户剩余积分:2

视频地址:http://wxapp.tc.qq.com/251/20302/stodownload?encfilekey=oibeqyX228riaCwo9STVsGLIBn9G5YG8ZnkvY5JFcntP76AqYZaFM8WXp3DrpxuP2DabEqsZAkHJu5gL31XVFyGdlO2zsXiamCOlgz0qOwg1rVKMbOkDYgf6dfMj8uH0nAX9fuhcPo1NeM&token=Cvvj5Ix3eezqxdzOWuG5sibF8S6CfvA5QdwjXV3mVBwAvVjwk8kkbdH8gwwVRAlF9QW0USm7o2C3uNqlKJPqwB1hVriaaDicv6vTJiaH4WvuVg8B5LrjvmpVhK6h2GDc1wB6&idx=1&a=1&bizid=1023&dotrans=0&hy=SZ&m=0cf4ff50ad161ac3ef9ddb1d13a5f751&upid=290280

封面地址:http://wxapp.tc.qq.com/251/20350/stodownload?encfilekey=oibeqyX228riaCwo9STVsGLIBn9G5YG8ZnqEKxaodOB3d7gxYw2psADjTia1sZMmufdLMUDgq782LygDKLPafmp3jobibDMajRodic9mpjLlII02zpzgQv7ibmoTib3a4DibPGyOZWGOliacLcSc&token=Cvvj5Ix3eeyD0TVgRZ2eE6Slcxmgc4IFWz0pChic8CLjOLA7EOW5NYzgqa8m5klKiap4vU6Ando7M13JOFcLJxyuj9icMfdibQtdLkJZR3kX2sUAe6iciaqDPsicr8oyZRnnsg1&idx=1&bizid=1023&dotrans=0&hy=SZ&m=904f9dc2936b2954287a5b1079a94c61

标题:这两片西瓜好像惹事了

主题:李沐-`
	resp := Extract(req, "视频地址:", "\n")
	t.Logf("resp:%v", resp)
}
