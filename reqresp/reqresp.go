package reqresp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type ErrMsg struct {
	Errcode int64  `json:"err_code"`
	Errmsg  string `json:"err_msg"`
}

func MakeResp(c *gin.Context, obj interface{}) {
	respJson, _ := json.Marshal(obj)
	c.Set(CtxFieldResp, string(respJson))
	c.JSON(200, obj)
}

func MakeRespOk(c *gin.Context) {
	MakeResp(c, struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	})
}

func MakeErrMsg(c *gin.Context, err error) {
	errObj := &ErrMsg{
		Errcode: -1,
		Errmsg:  err.Error(),
	}
	MakeResp(c, errObj)
}

func ParseReq(c *gin.Context, req interface{}) error {
	rawData, err := c.GetRawData()
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawData, req)
	if err != nil {
		return err
	}

	return nil
}
