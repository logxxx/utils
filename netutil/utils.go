package netutil

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Err404     = errors.New("404")
	Err302     = errors.New("302")
	ErrTimeout = errors.New("timeout")
	ErrNoHost  = errors.New("no host")
)

func GetImage(url string) (respBytes []byte, err error) {

	retryTimes := 0
RETRY:
	respBody, err := httpGet(nil, url)
	if err != nil {
		log.Printf("GetImage httpGet err:%v url:%v", err, url)
		if err == Err404 || err == ErrTimeout || err == ErrNoHost {
			retryTimes++
			if retryTimes < 3 {
				log.Printf("RETRY %v", retryTimes)
				goto RETRY
			}
		} else {
			return nil, err
		}
	}
	if respBody == nil {
		return nil, nil
	}
	defer respBody.Close()

	respBytes, err = ioutil.ReadAll(respBody)
	if err != nil {
		log.Printf("GetImage ReadAll err:%v url:%v", err, url)
		return nil, err
	}

	return respBytes, nil
}

func completePath(path string) string {
	return joinPath("http://q.quantuwang1.com", path)
}

func genCheckRedirectFunc(referer string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Referer", referer)
		return nil
	}
}

func joinPath(elem ...string) string {
	return strings.Join(elem, "")
}

func httpGet(client *http.Client, url string) (io.ReadCloser, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	httpResp, err := client.Do(req)
	if err != nil {
		log.Printf("doFind Get err:%v", err)
		if strings.Contains(err.Error(), "Timeout") {
			err = ErrTimeout
			return nil, err
		}
		if strings.Contains(err.Error(), "no such host") {
			err = ErrNoHost
			return nil, err
		}

		return nil, err
	}

	if httpResp.Status == "404" {
		err = Err404
		log.Printf("doFind Get err:%v", err)
		return nil, Err404
	}

	if httpResp.StatusCode == 302 {
		return nil, Err302
	}

	return httpResp.Body, nil
}

func SetHttpProxy(proxyURL string) (httpclient *http.Client) {
	ProxyURL, _ := url.Parse(proxyURL)
	httpclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
		},
	}
	return httpclient
}

func DownloadImage(url, path string) error {
	data, err := GetImage(url)
	if err != nil {
		return err
	}
	if len(data) < 1024*200 { //尺寸太小
		return nil
	}
	os.MkdirAll(filepath.Dir(path), 0777)
	return ioutil.WriteFile(path, data, 0777)
}
