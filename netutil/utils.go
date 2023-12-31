package netutil

import (
	"errors"
	"fmt"
	"io"
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

func httpGet(client *http.Client, url string, setHeaderFn func(r *http.Request)) (io.ReadCloser, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if setHeaderFn != nil {
		setHeaderFn(req)
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

func DownloadToFile(url string, downloadPath string, httpClient ...*http.Client) error {
	c := http.DefaultClient
	if len(httpClient) > 0 {
		c = httpClient[0]
	}
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid code:%v", resp.StatusCode)
	}

	os.MkdirAll(filepath.Dir(downloadPath), 0755)

	os.Remove(downloadPath)

	f, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
