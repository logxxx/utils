package log

import (
	"encoding/json"
	"log"
)

func Info(v ...interface{}) {
	log.Print(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func JSON(obj interface{}) string {
	resp, _ := json.Marshal(obj)
	return string(resp)
}

func JSONGrace(obj interface{}) string {
	resp, _ := json.MarshalIndent(obj, "", " ")
	return string(resp)
}

func Dumb(format string, v ...interface{}) {
	//DO NOTHING
}

func Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Error(v ...interface{}) {
	log.Print(v...)
}

func Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
