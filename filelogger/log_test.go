package filelogger_test

import (
	"github.com/logxxx/utils/filelogger"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger1, err := filelogger.NewLogger("test1")
	if err != nil {
		t.Fatal(err)
	}
	logger2, err := filelogger.NewLogger("test2")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		logger1.WithField("hello", i).Infof("im logger%v", 1)
		logger2.WithField("world", i).Infof("im logger%v", 2)
	}

	t.Logf("logger1:%v", logger1.GetFromFile())
	t.Logf("logger2:%v", logger2.GetFromFile())

	//logger1.Close()
	//logger2.Close()

}
