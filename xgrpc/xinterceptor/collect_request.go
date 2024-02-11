package xinterceptor

import (
	"context"
	"fmt"

	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xcommon"
	"github.com/blackRice-Tu/golib/xgrpc"
	"github.com/blackRice-Tu/golib/xlogger"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InterceptorLogger(l *zap.Logger) grpcLogging.Logger {
	return grpcLogging.LoggerFunc(func(ctx context.Context, lvl grpcLogging.Level, msg string, fields ...any) {
		body := map[string]any{
			"msg": msg,
		}
		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			body[key.(string)] = value
		}
		bodyJson := xcommon.JsonMarshal(body)

		switch lvl {
		case grpcLogging.LevelDebug:
			l.Debug(bodyJson)
		case grpcLogging.LevelInfo:
			l.Info(bodyJson)
		case grpcLogging.LevelWarn:
			l.Warn(bodyJson)
		case grpcLogging.LevelError:
			l.Error(bodyJson)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

type CollectRequestOpt struct {
	EscapeFunc     func(ctx context.Context) bool
	CustomInfoFunc func(ctx context.Context) map[string]any
	LoggerConfig   *xlogger.Config
}

func CollectRequestInterceptor(opt *CollectRequestOpt) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	if opt == nil {
		opt = &CollectRequestOpt{}
	}
	logger := InterceptorLogger(xlogger.NewLogger(opt.LoggerConfig))

	var interceptorOpts []grpcLogging.Option
	var interceptorOpt grpcLogging.Option

	// event
	interceptorOpt = grpcLogging.WithLogOnEvents(grpcLogging.StartCall, grpcLogging.FinishCall, grpcLogging.PayloadReceived, grpcLogging.PayloadSent)
	interceptorOpts = append(interceptorOpts, interceptorOpt)

	// timeFmt
	interceptorOpt = grpcLogging.WithTimestampFormat(golib.GetTimeFmt())
	interceptorOpts = append(interceptorOpts, interceptorOpt)

	// custom info
	customFields := func(ctx context.Context) grpcLogging.Fields {
		fields := grpcLogging.Fields{
			"env", golib.GetEnv(),
			"server_ip", golib.GetServerIp(),
			"program_key", golib.GetProgramKey(),
			"trace_id", xgrpc.ServerGetTraceId(ctx),
		}
		if opt.CustomInfoFunc != nil {
			fields.AppendUnique(grpcLogging.Fields{"custom", opt.CustomInfoFunc(ctx)})
		}
		return fields
	}
	interceptorOpt = grpcLogging.WithFieldsFromContext(customFields)
	interceptorOpts = append(interceptorOpts, interceptorOpt)

	ui := grpcLogging.UnaryServerInterceptor(logger, interceptorOpts...)
	si := grpcLogging.StreamServerInterceptor(logger, interceptorOpts...)
	return ui, si
}
