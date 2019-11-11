package webutil

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
)

// NewGracefulServer returns a new graceful server.
func NewGracefulServer(server *http.Server) *GracefulServer {
	return &GracefulServer{
		latch:  &async.Latch{},
		server: server,
	}
}

// GracefulServer is a wrapper for an http server that implements the graceful interface.
type GracefulServer struct {
	latch               *async.Latch
	server              *http.Server
	shutdownGracePeriod time.Duration
	listener            net.Listener
}

// WithShutdownGracePeriod sets the shutdown grace period.
func (gs *GracefulServer) WithShutdownGracePeriod(d time.Duration) *GracefulServer {
	gs.shutdownGracePeriod = d
	return gs
}

// ShutdownGracePeriod returns the shutdown graceperiod or a default.
func (gs *GracefulServer) ShutdownGracePeriod() time.Duration {
	return gs.shutdownGracePeriod
}

// WithListener sets the server listener.
func (gs *GracefulServer) WithListener(l net.Listener) *GracefulServer {
	gs.listener = l
	return gs
}

// Listener returns the listener.
func (gs *GracefulServer) Listener() net.Listener {
	return gs.listener
}

// Start implements graceful.Graceful.Start.
// It is expected to block.
func (gs *GracefulServer) Start() (err error) {
	gs.latch.Started()

	var shutdownErr error
	if gs.listener != nil {
		shutdownErr = gs.server.Serve(gs.listener)
	} else {
		shutdownErr = gs.server.ListenAndServe()
	}

	gs.latch.Stopped()
	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		err = exception.New(shutdownErr)
	}
	return
}

// Stop implements graceful.Graceful.Stop.
func (gs *GracefulServer) Stop() error {
	if !gs.latch.IsRunning() {
		return nil
	}
	gs.latch.Stopping()
	gs.server.SetKeepAlivesEnabled(false)

	ctx := context.Background()
	if gs.shutdownGracePeriod > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, gs.shutdownGracePeriod)
		defer cancel()
	}
	return exception.New(gs.server.Shutdown(ctx))
}

// NotifyStarted implements graceful.Graceful.NotifyStarted.
func (gs *GracefulServer) NotifyStarted() <-chan struct{} {
	return gs.latch.NotifyStarted()
}

// NotifyStopped implements graceful.Graceful.NotifyStopped.
func (gs *GracefulServer) NotifyStopped() <-chan struct{} {
	return gs.latch.NotifyStopped()
}
