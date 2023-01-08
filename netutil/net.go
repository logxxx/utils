package netutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpDo(req *http.Request) (int, []byte, error) {
	httpResp, err := http.DefaultClient.Do(req)
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

func HttpGetRaw(url string) (int, []byte, error) {
	httpResp, err := http.Get(url)
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

func HttpReqGet(req *http.Request, resp interface{}) (int, error) {
	status, respBytes, err := HttpDo(req)
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

func HttpGet(url string, resp interface{}) (int, error) {

	status, respBytes, err := HttpGetRaw(url)
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

func HttpPost(url string, reqBody interface{}, resp interface{}) (int, error) {

	reqBodyBytes := make([]byte, 0)
	var err error
	if reqBody != nil {
		reqBodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			return 0, err
		}
	}

	reqBodyBuf := bytes.NewBuffer(reqBodyBytes)
	httpResp, err := http.Post(url, "application/json", reqBodyBuf)
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
