package signal

import (
	"context"
	"os"
	osig "os/signal"
	"syscall"
)

func WaitForInterrupt(cb func()) {
	ch := make(chan os.Signal)
	osig.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	if cb != nil {
		cb()
	}
}

func WithSignal(parent context.Context) context.Context {
	return WithSignalEx(parent, nil)
}

func WithSignalEx(parent context.Context, cb func()) context.Context {
	ch := make(chan os.Signal)
	osig.Notify(ch, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-ch:
			if cb != nil {
				cb()
			}
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
