package xcontext

import (
	"context"

	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xalgorithm"
)

func GenerateTraceId() string {
	return xalgorithm.GetRandString(32)
}

func NewTraceContext() (context.Context, string) {
	traceId := GenerateTraceId()
	ctx := context.WithValue(context.Background(), golib.GetTraceIdKey(), traceId)
	return ctx, traceId
}

func NewTraceContextWithCancel() (context.Context, string, context.CancelFunc) {
	parentCtx, traceId := NewTraceContext()
	ctx, cancel := context.WithCancel(parentCtx)
	return ctx, traceId, cancel
}

func GetOrNewTraceId(ctx context.Context) string {
	traceId := ""
	traceIdKey := golib.GetTraceIdKey()
	traceIdValue := ctx.Value(traceIdKey)
	if traceIdValue != nil {
		traceId = traceIdValue.(string)
	}
	if traceId == "" {
		traceId = GenerateTraceId()
	}
	return traceId
}

func SetTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, golib.GetTraceIdKey(), traceId)
}
