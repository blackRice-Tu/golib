package xgrpc

import (
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	_ "google.golang.org/grpc/metadata"
)

// Dialer for vrpc client
type Dialer struct {
	target                   string
	opts                     []grpc.DialOption
	unaryClientInterceptors  []grpc.UnaryClientInterceptor
	streamClientInterceptors []grpc.StreamClientInterceptor
}

type DialerRegistry func(d *Dialer)

// NewDialer xx
func NewDialer(dr DialerRegistry) *Dialer {
	d := &Dialer{
		opts:                     []grpc.DialOption{},
		unaryClientInterceptors:  make([]grpc.UnaryClientInterceptor, 0),
		streamClientInterceptors: make([]grpc.StreamClientInterceptor, 0),
	}
	dr(d)
	return d
}

// DirectRegistry xx
func DirectRegistry(target string) DialerRegistry {
	return func(d *Dialer) {
		d.target = target
	}
}

// UseUnaryInterceptor xx
func (d *Dialer) UseUnaryInterceptor(interceptors ...grpc.UnaryClientInterceptor) {
	d.unaryClientInterceptors = append(d.unaryClientInterceptors, interceptors...)
}

// UseStreamInterceptor xx
func (d *Dialer) UseStreamInterceptor(interceptors ...grpc.StreamClientInterceptor) {
	d.streamClientInterceptors = append(d.streamClientInterceptors, interceptors...)
}

// WithOption xx
func (d *Dialer) WithOption(opt ...grpc.DialOption) {
	d.opts = append(d.opts, opt...)
}

// Dial xx
func (d *Dialer) Dial() (*grpc.ClientConn, error) {
	// keepalive
	keepaliveOpt := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                150 * time.Second,
		Timeout:             250 * time.Millisecond,
		PermitWithoutStream: true,
	})
	chainUnaryInterceptor := middleware.ChainUnaryClient(d.unaryClientInterceptors...)
	chainStreamInterceptor := middleware.ChainStreamClient(d.streamClientInterceptors...)
	d.WithOption(
		keepaliveOpt,
		grpc.WithUnaryInterceptor(chainUnaryInterceptor),
		grpc.WithStreamInterceptor(chainStreamInterceptor))
	return grpc.Dial(d.target, d.opts...)
}
