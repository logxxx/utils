package runutil

import (
	"fmt"
	"github.com/logxxx/utils/log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func GoRunSafe(fn func()) {
	if fn == nil {
		return
	}
	func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				_ = runtime.Stack(buf, false)
				log.Errorf("RunSafe panic:%v stack:%s", err, buf)
			}
		}()

		fn()
	}()
}

// 安全地执行goroutine。包装了recover()捕获异常
func RunSafe(fn func()) {
	if fn == nil {
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				_ = runtime.Stack(buf, false)
				log.Errorf("RunSafe panic:%v stack:%s", err, buf)
			}
		}()

		fn()
	}()
}

func WaitForExit(closeFunc func()) {
	doneChan := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-signalChan:
			fmt.Printf("captured %v. exiting...\n", s)
			if closeFunc != nil {
				closeFunc()
			}
			close(doneChan)
		case <-doneChan:
			os.Exit(0)
		}
	}
}
