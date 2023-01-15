package netutil

import (
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/reqresp"
	"sort"
	"time"
)

type IPLogger struct {
	Statistic Statistic
	Records   []*Record
}

type Statistic struct {
	CreateTime     int64
	CreateTimeStr  string
	TotalCount     int64
	CountPerMinute map[int64]int64
	CountPerHour   map[int64]int64
}

type Record struct {
	CreateTime    int64
	CreateTimeStr string
	ReqIP         string
	ReqURL        string
	ReqURI        string
}

type GetRecordsResp struct {
	Records []*Record
}

type GetStatisticResp struct {
	CreateTimeStr   string
	TotalCount      int64
	CountsPerMinute []*CountPerTimeSlice
	CountsPerHour   []*CountPerTimeSlice
}

type CountPerTimeSlice struct {
	Time    int64
	TimeStr string
	Count   int64
}

var (
	_ipLogger *IPLogger
)

func NewIPLogger() *IPLogger {
	l := &IPLogger{
		Statistic: Statistic{
			CreateTime:     time.Now().Unix(),
			CreateTimeStr:  time.Now().Format("2006-01-02 15:04:05"),
			TotalCount:     0,
			CountPerMinute: make(map[int64]int64, 0),
			CountPerHour:   make(map[int64]int64, 0),
		},
		Records: nil,
	}
	return l
}

func GetIPLogger() *IPLogger {
	if _ipLogger != nil {
		return _ipLogger
	}
	_ipLogger = NewIPLogger()
	return _ipLogger
}

func (s *Statistic) Add() {
	s.TotalCount++
	key := time.Now().Unix() / 60 * 60
	s.CountPerMinute[key]++
	key = time.Now().Unix() / 3600 * 3600
	s.CountPerHour[key]++
}

func (l *IPLogger) Log(c *gin.Context) {
	reqIP := c.RemoteIP()
	reqURL := c.Request.URL.String()
	reqURI := c.Request.RequestURI
	now := time.Now()
	record := &Record{
		CreateTime:    now.Unix(),
		CreateTimeStr: now.Format("2006-01-02 15:04:05"),
		ReqIP:         reqIP,
		ReqURL:        reqURL,
		ReqURI:        reqURI,
	}

	l.Records = append([]*Record{record}, l.Records...)
	l.Statistic.Add()
}

func (l *IPLogger) RegisterAPI_Clean(c *gin.Context) {
	l.Clean()
	reqresp.MakeRespOk(c)
}

func (l *IPLogger) RegisterAPI_GetRecords(c *gin.Context) {
	resp := &GetRecordsResp{
		Records: l.Records,
	}
	reqresp.MakeResp(c, resp)
}

func (l *IPLogger) RegisterAPI_GetStatistic(c *gin.Context) {
	resp := &GetStatisticResp{
		CreateTimeStr: l.Statistic.CreateTimeStr,
		TotalCount:    l.Statistic.TotalCount,
	}

	countsPerMinute := make([]*CountPerTimeSlice, 0)
	for min, count := range l.Statistic.CountPerMinute {
		countsPerMinute = append(countsPerMinute, &CountPerTimeSlice{
			Time:    min,
			TimeStr: time.Unix(min, 0).Format("01/02 15:04"),
			Count:   count,
		})
	}

	sort.Slice(countsPerMinute, func(i, j int) bool { //倒序
		return countsPerMinute[i].Time > countsPerMinute[j].Time
	})

	resp.CountsPerMinute = countsPerMinute

	countsPerHour := make([]*CountPerTimeSlice, 0)
	for h, count := range l.Statistic.CountPerHour {
		countsPerHour = append(countsPerHour, &CountPerTimeSlice{
			Time:    h,
			TimeStr: time.Unix(h, 0).Format("01/02 15:04"),
			Count:   count,
		})
	}

	sort.Slice(countsPerHour, func(i, j int) bool { //倒序
		return countsPerHour[i].Time > countsPerHour[j].Time
	})

	resp.CountsPerHour = countsPerHour

	reqresp.MakeResp(c, resp)

}

func (l *IPLogger) Clean() {
	l = NewIPLogger()
}
