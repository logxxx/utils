package netutil

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	respBody, err := httpGet(url, true)
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

func DoFind(url string, findFunc func(doc *goquery.Document) error) (err error) {

	respBody, err := httpGet(url, false)
	if err != nil {
		return err
	}

	defer respBody.Close()

	//utf8Reader := transform.NewReader(respBody, simplifiedchinese.GBK.NewDecoder())

	doc, err := goquery.NewDocumentFromReader(respBody)
	if err != nil {
		log.Printf("doFind NewDocumentFromReader err:%v", err)
		if strings.Contains(err.Error(), "Timeout") {
			err = ErrTimeout
			return err
		}
		return err
	}

	err = findFunc(doc)
	if err != nil {
		log.Printf("doFind findFunc err:%v", err)
		return err
	}

	return nil
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

func httpGet(url string, isImage bool) (io.ReadCloser, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 10 * time.Second}
	if isImage {
		req.Header.Set("Authority", "p.xgmn.vip")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Cookie", "Hm_lvt_3e4b7c3cd2459ed4a577e4795c1973f9=1660411375,1660412551,1660416864,1660971543; Hm_lpvt_3e4b7c3cd2459ed4a577e4795c1973f9=1660972432")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	} else {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Cookie", "Hm_lvt_3e4b7c3cd2459ed4a577e4795c1973f9=1660058145,1660379072,1660411375; Hm_lpvt_3e4b7c3cd2459ed4a577e4795c1973f9=1660411425")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
		req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
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

type DocFinder struct {
	cache map[string]*goquery.Document
}

func NewDocFinder() *DocFinder {
	return &DocFinder{
		cache: make(map[string]*goquery.Document),
	}
}

func (f *DocFinder) Find(url string, findFunc func(doc *goquery.Document) error) (err error) {

	var document *goquery.Document

	if value, ok := f.cache[url]; ok {
		document = value
	} else {
		respBody, err := httpGet(url, false)
		if err != nil {
			return err
		}

		defer respBody.Close()

		//utf8Reader := transform.NewReader(respBody, simplifiedchinese.GBK.NewDecoder())

		doc, err := goquery.NewDocumentFromReader(respBody)
		if err != nil {
			log.Printf("doFind NewDocumentFromReader err:%v", err)
			if strings.Contains(err.Error(), "Timeout") {
				err = ErrTimeout
				return err
			}
			return err
		}
		document = doc
	}

	err = findFunc(document)
	if err != nil {
		log.Printf("doFind findFunc err:%v", err)
		return err
	}

	return nil

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
