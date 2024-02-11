package xinterceptor

import (
	"context"
	"github.com/blackRice-Tu/golib/utils/xcontext"
	"github.com/blackRice-Tu/golib/xgrpc"
	logger "github.com/blackRice-Tu/golib/xlogger/default"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func RecoveryInterceptor() (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	var interceptorOpts []grpcRecovery.Option
	var interceptorOpt grpcRecovery.Option

	// recovery handler
	handler := func(ctx context.Context, p any) (err error) {
		ctx = xcontext.SetTraceId(ctx, xgrpc.ServerGetTraceId(ctx))
		logger.Errorf(ctx, "gRPC server error: %v\n%s", p, string(debug.Stack()))
		return status.Errorf(codes.Unknown, "system error: %v", p)
	}
	interceptorOpt = grpcRecovery.WithRecoveryHandlerContext(handler)
	interceptorOpts = append(interceptorOpts, interceptorOpt)

	ui := grpcRecovery.UnaryServerInterceptor(interceptorOpts...)
	si := grpcRecovery.StreamServerInterceptor(interceptorOpts...)
	return ui, si
}
