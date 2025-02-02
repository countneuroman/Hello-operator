package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var signalHandler = make(chan struct{}) //minimal channel for synchronization

func SetupSignalHandler() context.Context {
	close(signalHandler)

	c := make(chan os.Signal, 2)
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	return ctx
}
