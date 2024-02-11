package xsys

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitSignal awaits for SIGINT or SIGTERM and closes the channel
func WaitSignal(stop chan struct{}) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	close(stop)
}
