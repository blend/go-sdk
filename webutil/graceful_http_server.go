/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// NewGracefulHTTPServer returns a new graceful http server wrapper.
func NewGracefulHTTPServer(server *http.Server, options ...GracefulHTTPServerOption) *GracefulHTTPServer {
	gs := &GracefulHTTPServer{
		Latch:  async.NewLatch(),
		Server: server,
	}
	for _, option := range options {
		option(gs)
	}
	return gs
}

// GracefulHTTPServerOption is an option for the graceful http server.
type GracefulHTTPServerOption func(*GracefulHTTPServer)

// OptGracefulHTTPServerShutdownGracePeriod sets the shutdown grace period.
func OptGracefulHTTPServerShutdownGracePeriod(d time.Duration) GracefulHTTPServerOption {
	return func(g *GracefulHTTPServer) { g.ShutdownGracePeriod = d }
}

// OptGracefulHTTPServerListener sets the server listener.
func OptGracefulHTTPServerListener(listener net.Listener) GracefulHTTPServerOption {
	return func(g *GracefulHTTPServer) { g.Listener = listener }
}

// OptGracefulHTTPServerLog sets the logger.
func OptGracefulHTTPServerLog(log logger.Log) GracefulHTTPServerOption {
	return func(g *GracefulHTTPServer) { g.Log = log }
}

// GracefulHTTPServer is a wrapper for an http server that implements the graceful interface.
type GracefulHTTPServer struct {
	Log                 logger.Log
	Latch               *async.Latch
	Server              *http.Server
	ShutdownGracePeriod time.Duration
	Listener            net.Listener
}

// Start implements graceful.Graceful.Start.
// It is expected to block.
func (gs *GracefulHTTPServer) Start() (err error) {
	if !gs.Latch.CanStart() {
		err = ex.New(async.ErrCannotStart)
		return
	}
	gs.Latch.Starting()
	gs.Latch.Started()
	defer gs.Latch.Stopped()

	var shutdownErr error
	if gs.Listener != nil {
		logger.MaybeInfof(gs.Log, "http server listening on %s", gs.Listener.Addr().String())
		shutdownErr = gs.Server.Serve(gs.Listener)
	} else {
		logger.MaybeInfof(gs.Log, "http server listening on %s", gs.Server.Addr)
		shutdownErr = gs.Server.ListenAndServe()
	}
	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		err = ex.New(shutdownErr)
	}
	return
}

// Stop implements graceful.Graceful.Stop.
func (gs *GracefulHTTPServer) Stop() error {
	if !gs.Latch.CanStop() {
		return ex.New(async.ErrCannotStop)
	}
	gs.Latch.Stopping()
	gs.Server.SetKeepAlivesEnabled(false)
	ctx := context.Background()
	if gs.ShutdownGracePeriod > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, gs.ShutdownGracePeriod)
		defer cancel()
	}
	return ex.New(gs.Server.Shutdown(ctx))
}

// NotifyStarted implements part of graceful.
func (gs *GracefulHTTPServer) NotifyStarted() <-chan struct{} {
	return gs.Latch.NotifyStarted()
}

// NotifyStopped implements part of graceful.
func (gs *GracefulHTTPServer) NotifyStopped() <-chan struct{} {
	return gs.Latch.NotifyStopped()
}
