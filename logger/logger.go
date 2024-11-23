package logger

import (
	"fmt"
	"github.com/logxxx/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sort"
	"strings"
)

var (
	pid = os.Getpid()
)

type MyLogFormatter struct {
}

func (f *MyLogFormatter) Format(entry *log.Entry) ([]byte, error) {

	level := "I"
	if entry.Level == log.DebugLevel {
		level = "D"
	} else if entry.Level == log.ErrorLevel {
		level = "E"
	} else if entry.Level == log.WarnLevel {
		level = "W"
	}

	datas := []string{}
	for k, v := range entry.Data {
		showV := fmt.Sprintf("%v", v)
		vKind := reflect.TypeOf(v).Kind()
		if vKind == reflect.Struct || vKind == reflect.Ptr {
			showV = utils.JsonToString(v)
		}
		datas = append(datas, fmt.Sprintf("%v=%v", k, showV))
	}
	sort.Strings(datas)

	resp := fmt.Sprintf("%v[%v]%v|%v %v\n", entry.Time.Format("01/02 15:04:05"), level, pid, strings.Join(datas, "&"), entry.Message)

	fmt.Print(resp)

	return []byte(resp), nil
}
