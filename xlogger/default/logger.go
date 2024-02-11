package _default

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/blackRice-Tu/golib/utils/xcontext"
	"github.com/blackRice-Tu/golib/xlogger"

	"go.uber.org/zap"
)

var (
	initOnce sync.Once
	logger   *zap.Logger
)

// InitLogger ...
func InitLogger(logConf *xlogger.Config) {
	initOnce.Do(func() {
		// zap json logger
		logger = xlogger.NewLogger(logConf)
	})
}

// GetLogger ...
func GetLogger() *zap.Logger {
	// redirect to stdout
	if logger == nil {
		stdoutLogger, _ := zap.NewProductionConfig().Build()
		return stdoutLogger
	}
	return logger
}

func formatMsg(ctx context.Context, content string) (m string) {
	traceId := xcontext.GetOrNewTraceId(ctx)
	msgField := struct {
		xlogger.CommonField
		Content string `json:"content"`
	}{
		CommonField: xlogger.GetCommonField(traceId),
		Content:     content,
	}
	msgBytes, _ := json.Marshal(msgField)
	return string(msgBytes)
}

// Debug ...
func Debug(ctx context.Context, a ...any) {
	GetLogger().Debug(formatMsg(ctx, fmt.Sprint(a...)))
}

// Debugf ...
func Debugf(ctx context.Context, msg string, a ...any) {
	GetLogger().Debug(formatMsg(ctx, fmt.Sprintf(msg, a...)))
}

// Info ...
func Info(ctx context.Context, a ...any) {
	GetLogger().Info(formatMsg(ctx, fmt.Sprint(a...)))
}

// Infof ...
func Infof(ctx context.Context, msg string, a ...any) {
	GetLogger().Info(formatMsg(ctx, fmt.Sprintf(msg, a...)))
}

// Warn ....
func Warn(ctx context.Context, a ...any) {
	GetLogger().Warn(formatMsg(ctx, fmt.Sprint(a...)))
}

// Warnf ...
func Warnf(ctx context.Context, msg string, a ...any) {
	GetLogger().Warn(formatMsg(ctx, fmt.Sprintf(msg, a...)))
}

// Error ...
func Error(ctx context.Context, a ...any) {
	GetLogger().Error(formatMsg(ctx, fmt.Sprint(a...)))
}

// Errorf ...
func Errorf(ctx context.Context, msg string, a ...any) {
	GetLogger().Error(formatMsg(ctx, fmt.Sprintf(msg, a...)))
}

// Sync ...
func Sync() error {
	return GetLogger().Sync()
}
