package grpcutil

import (
	"net"

	"google.golang.org/grpc"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/logger"
)

// NewGraceful returns a new graceful host for a grpc server.
func NewGraceful(listener net.Listener, server *grpc.Server) *Graceful {
	return &Graceful{
		Latch:    async.NewLatch(),
		Listener: listener,
		Server:   server,
	}
}

// Graceful is a shim for graceful hosting grpc servers.
type Graceful struct {
	*async.Latch
	Log      logger.Log
	Listener net.Listener
	Server   *grpc.Server
}

// WithLogger sets the logger.
func (gz *Graceful) WithLogger(log logger.Log) *Graceful {
	gz.Log = log
	return gz
}

// Start starts the server.
func (gz *Graceful) Start() error {
	gz.Latch.Starting()
	gz.Latch.Started()
	logger.MaybeInfof(gz.Log, "grpc server starting, listening on %v %s", gz.Listener.Addr().Network(), gz.Listener.Addr().String())
	return gz.Server.Serve(gz.Listener)
}

// Stop shuts the server down.
func (gz *Graceful) Stop() error {
	gz.Latch.Stopping()
	logger.MaybeInfof(gz.Log, "grpc server shutting down")
	gz.Server.GracefulStop()
	gz.Latch.Stopped()
	return nil
}

// IsRunning returns if the server is running.
func (gz *Graceful) IsRunning() bool {
	return gz.IsRunning()
}

// NotifyStarted returns the notify started signal.
func (gz *Graceful) NotifyStarted() <-chan struct{} {
	return gz.NotifyStarted()
}

// NotifyStopped returns the notify stopped signal.
func (gz *Graceful) NotifyStopped() <-chan struct{} {
	return gz.NotifyStopped()
}
