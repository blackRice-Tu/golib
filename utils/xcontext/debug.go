package xcontext

import (
	"context"

	"github.com/blackRice-Tu/golib"
)

const (
	debugFlag = "1"
)

func NewDebugContext() context.Context {
	return context.WithValue(context.Background(), golib.GetDebugKey(), debugFlag)
}

func SetDebug(ctx context.Context) context.Context {
	return context.WithValue(ctx, golib.GetDebugKey(), debugFlag)
}

func IsDebug(ctx context.Context) (ok bool) {
	debugKey := golib.GetDebugKey()
	debugValue := ctx.Value(debugKey)
	if debugValue != nil && debugValue.(string) == debugFlag {
		ok = true
	}
	return
}
