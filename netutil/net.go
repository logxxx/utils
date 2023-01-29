package netutil

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

func TryGetFileSize(url string, setHeaderFuncs ...func(httpReq *http.Request)) (fileSize int64) {

	httpReq, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}

	if len(setHeaderFuncs) > 0 {
		setHeaderFuncs[0](httpReq)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Errorf("TryGetFileSize err:%v", err)
		return
	}

	log.Infof("TryGetFileSize resp.Code:%v", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return
	}

	contentRangeStr := resp.Header.Get("Content-Length")

	size, _ := strconv.ParseInt(contentRangeStr, 10, 64)

	return size

}

func HttpDo(req *http.Request, httpClient ...*http.Client) (int, []byte, error) {
	client := http.DefaultClient
	if len(httpClient) > 0 {
		client = httpClient[0]
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}

	return httpResp.StatusCode, respBytes, nil

}

func HttpGetRaw(url string, httpClient ...*http.Client) (int, []byte, error) {

	client := http.DefaultClient
	if len(httpClient) > 0 {
		client = httpClient[0]
	}

	httpResp, err := client.Get(url)
	if err != nil {
		return 0, nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return 0, nil, err
	}

	return httpResp.StatusCode, respBytes, nil
}

func HttpReqGet(req *http.Request, resp interface{}, httpClient ...*http.Client) (int, error) {

	client := http.DefaultClient
	if len(httpClient) > 0 {
		client = httpClient[0]
	}

	status, respBytes, err := HttpDo(req, client)
	if err != nil {
		return 0, err
	}

	if resp != nil {
		err = json.Unmarshal(respBytes, resp)
		if err != nil {
			return 0, err
		}
	}

	return status, nil

}

func HttpGet(url string, resp interface{}, httpClient ...*http.Client) (int, error) {

	client := http.DefaultClient
	if len(httpClient) > 0 {
		client = httpClient[0]
	}

	status, respBytes, err := HttpGetRaw(url, client)
	if err != nil {
		return 0, err
	}

	//log.Debugf("HttpGet\nreq:%v\nresp:%v", url, string(respBytes))

	if resp != nil {
		err = json.Unmarshal(respBytes, resp)
		if err != nil {
			return 0, err
		}
	}

	return status, nil

}

func HttpPost(url string, reqBody interface{}, resp interface{}, httpClient ...*http.Client) (int, error) {

	client := http.DefaultClient
	if len(httpClient) > 0 {
		client = httpClient[0]
	}

	reqBodyBytes := make([]byte, 0)
	var err error
	if reqBody != nil {
		reqBodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			return 0, err
		}
	}

	reqBodyBuf := bytes.NewBuffer(reqBodyBytes)
	httpResp, err := client.Post(url, "application/json", reqBodyBuf)
	if err != nil {
		return 0, err
	}
	defer httpResp.Body.Close()

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return 0, err
	}

	if resp != nil {
		err = json.Unmarshal(respBytes, resp)
		if err != nil {
			return 0, err
		}
	}

	return httpResp.StatusCode, nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requ ested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
