package filelogger

import "log"

var (
	_default *Logger
)

func initDefaultLogger() {
	_default, _ = NewLogger("default")
}

func Infof(format string, args ...interface{}) {

	if _default == nil {
		initDefaultLogger()
	}

	prefix := "[info]"
	log.Printf(prefix+format, args...)
	if _default != nil {
		_default.writeToFile(prefix+format, args...)
	}
}
