package logutil

import (
	"context"
	"fmt"
	"github.com/logxxx/utils/randutil"
	log "github.com/sirupsen/logrus"
	"time"
)

func Log(funcName string) *log.Entry {

	return log.WithField("func_name", funcName)
}

func CtxLog(ctx context.Context, funcName string) (*log.Entry, context.Context) {

	if ctx == nil {
		ctx = context.Background()
	}

	traceID, ok := ctx.Value("trace_id").(string)
	if !ok {
		traceID = randutil.RandStr(8)
		ctx = context.WithValue(ctx, "trace_id", traceID)
	}
	logger := log.WithField("trace_id", fmt.Sprintf("%v", traceID)).WithField("func_name", funcName).WithField("func_st", time.Now().Format("15:04:05"))
	return logger, ctx
}
