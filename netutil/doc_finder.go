package netutil

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type DocFinder struct {
	headerFn func(*http.Request)
	proxy    *http.Client
	cache    map[string]*goquery.Document
}

func NewDocFinder() *DocFinder {
	return &DocFinder{
		cache: make(map[string]*goquery.Document),
	}
}

func (f *DocFinder) SetProxy(proxyURL string) *DocFinder {
	f.proxy = SetHttpProxy(proxyURL)
	return f
}

func (f *DocFinder) SetHeader(fn func(req *http.Request)) {
	f.headerFn = fn
}

func (f *DocFinder) Find(url string, findFunc func(doc *goquery.Document) error) (err error) {

	var document *goquery.Document

	if value, ok := f.cache[url]; ok {
		document = value
	} else {
		respBody, err := httpGet(f.proxy, url, f.headerFn)
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
