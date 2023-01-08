package filelogger

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	loggerMap   = make(map[string]*Logger, 0)       //k:bizName
	logFileLock = make(map[string]*sync.RWMutex, 0) //k:fileName
)

type Logger struct {
	bizName  string
	fileName string
	file     *os.File
	fields   []Field
}

type Field struct {
	Key   string
	Value interface{}
}

func (f *Field) Format() string {
	result := ""
	if f.Value != "" {
		result = fmt.Sprintf("[%v=%v]", f.Key, f.Value)
	} else {
		result = fmt.Sprintf("[%v]", f.Key)
	}
	return result
}

func NewLogger(bizName string) (*Logger, error) {
	logger := loggerMap[bizName]
	if logger != nil {
		return logger, nil
	}

	name := getLoggerName(bizName)
	os.Remove(name)
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	logger = &Logger{
		bizName:  bizName,
		fileName: name,
		file:     f,
	}

	loggerMap[bizName] = logger
	logFileLock[name] = &sync.RWMutex{}

	return logger, nil

}

func (l *Logger) WithField(key string, values ...interface{}) *Logger {

	if key == "" {
		return l
	}

	for _, field := range l.fields {
		if key == field.Key {
			return l
		}
	}

	newField := Field{
		Key: key,
	}

	if len(values) > 0 {
		newField.Value = values[0]
	}

	newLogger := &Logger{
		bizName:  l.bizName,
		fileName: l.fileName,
		file:     l.file,
		fields:   append(l.fields, newField),
	}
	return newLogger
}

func (l *Logger) Infof(format string, args ...interface{}) {
	prefix := "[info]"
	for _, field := range l.fields {
		prefix += field.Format()
	}
	log.Infof(prefix+format, args...)
	l.writeToFile(prefix+format, args...)
}

func (l *Logger) writeToFile(format string, args ...interface{}) {

	logFileLock[l.fileName].Lock()
	defer logFileLock[l.fileName].Unlock()

	if l.file == nil {
		return
	}

	if !fileutil.HasFile(l.fileName) {
		return
	}

	content := fmt.Sprintf(format, args...)
	now := time.Now().Format("[15:04:05]")
	_, _ = l.file.WriteString(now + content + "\n\n")

}

func (l *Logger) GetFromFile() string {

	logFileLock[l.fileName].RLock()
	defer logFileLock[l.fileName].RUnlock()

	content, _ := ioutil.ReadFile(l.fileName)
	return string(content)

}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}

	os.Remove(l.fileName)

	delete(loggerMap, l.bizName)
}

func (l *Logger) Clean() error {

	logFileLock[l.fileName].Lock()
	defer logFileLock[l.fileName].Unlock()

	l.Close()

	newL, err := NewLogger(l.bizName)
	if err != nil {
		return err
	}

	l = newL

	return nil

}

func getLoggerName(bizName string) string {
	name := fmt.Sprintf("filelog_%v_%v.txt", bizName, time.Now().Format("20060102_150405"))
	return name
}

func getLogger(bizName string) *Logger {
	return loggerMap[bizName]
}
