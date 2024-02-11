package xgrpc

import (
	"context"

	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xcontext"

	"google.golang.org/grpc/metadata"
)

func WithTraceId(ctx context.Context) context.Context {
	traceId := xcontext.GetOrNewTraceId(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, golib.GetTraceIdKey(), traceId)
	return ctx
}

func ClientSetTraceId(ctx context.Context) context.Context {
	traceId := xcontext.GetOrNewTraceId(ctx)
	ctx = metadata.AppendToOutgoingContext(ctx, golib.GetTraceIdKey(), traceId)
	return ctx
}

func ServerGetTraceId(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	traceIdList := md.Get(golib.GetTraceIdKey())
	if len(traceIdList) == 0 {
		return xcontext.GenerateTraceId()
	}
	return traceIdList[0]
}
