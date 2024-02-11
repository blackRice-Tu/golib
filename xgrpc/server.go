package xgrpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xcommon"
	logger "github.com/blackRice-Tu/golib/xlogger/default"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// ServerConfig ...
type ServerConfig struct {
	RpcPort int
	Host    string
}

// RunServer ...
func RunServer(ctx context.Context, server *grpc.Server, sc *ServerConfig) error {
	return RunServers(ctx, []*grpc.Server{server}, []*ServerConfig{sc})
}

// RunServers ...
func RunServers(ctx context.Context, servers []*grpc.Server, scs []*ServerConfig) error {
	if len(scs) != len(servers) {
		return errors.Errorf("servers and configs num mismatched")
	}

	var wg sync.WaitGroup
	for i, server := range servers {
		sc := scs[i]
		addr := fmt.Sprintf("%s:%d", sc.Host, sc.RpcPort)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Errorf(ctx, "NetListenError: %v, addr: %s", err, addr)
			continue
		}
		listServices(ctx, server, sc)

		wg.Add(1)
		go func(i int, serve *grpc.Server, address string) {
			err := serve.Serve(l)
			if err != nil {
				logger.Errorf(ctx, "NetServeError: %v, addr: %s", err, addr)
			}
			wg.Done()
		}(i, server, addr)
	}
	wg.Wait()
	return nil
}

func listServices(ctx context.Context, server *grpc.Server, sc *ServerConfig) {
	serviceInfoMap := server.GetServiceInfo()
	if golib.GetMode() == golib.DebugMode {
		totalServiceNum := 0
		totalMethodsNum := 0
		title := fmt.Sprintf("[%s][%s:%d]", "GRPC-debug", sc.Host, sc.RpcPort)
		fmt.Printf("\n%s\t\n", title)
		for key, serviceInfo := range serviceInfoMap {
			totalServiceNum += 1
			for _, method := range serviceInfo.Methods {
				totalMethodsNum += 1
				fmt.Printf("%s\t%s\t%s\t%s\n", title, key, method.Name, xcommon.JsonMarshal(method))
			}
		}
		fmt.Printf("%s\tTotalServicesNum: %d\tTotalMethodsNum: %d\n", title, totalServiceNum, totalMethodsNum)
	}
}
