package reqresp

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"time"
)

const (
	ShowRespMaxLen = 512
	CtxFieldResp   = "rawRespJson"
)

func PrintReq(c *gin.Context) {
	st := time.Now()

	rawReq, _ := c.GetRawData()

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(rawReq))

	c.Next()

	showResp := ""
	respObj, ok := c.Get(CtxFieldResp)
	if ok {
		showResp = fmt.Sprintf("%v", respObj)
		if len(showResp) > ShowRespMaxLen {
			showResp = showResp[:ShowRespMaxLen] + fmt.Sprintf("...(%v more words ignored)", len(showResp)-ShowRespMaxLen)
		}
	}

	d := time.Since(st)

	log.Infof("[Req]url:%v st=%.2fs rawReq:%v rawResp:%v",
		c.Request.URL.String(), d.Seconds(), string(rawReq), showResp)
}
