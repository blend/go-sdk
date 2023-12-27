/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	radix "github.com/mediocregopher/radix/v4"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
)

var (
	_	async.Checker	= (*RadixClient)(nil)
	_	Client		= (*RadixClient)(nil)
)

// New returns a new client.
func New(ctx context.Context, opts ...Option) (*RadixClient, error) {
	var rc RadixClient
	var err error
	for _, opt := range opts {
		if err = opt(&rc); err != nil {
			return nil, ex.New(err)
		}
	}
	if rc.Config.ConnectTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, rc.Config.ConnectTimeout)
		defer cancel()
	}

	var dialer RadixNetDialer
	if rc.Config.UseTLS {
		dialer = new(tls.Dialer)
	}

	if len(rc.Config.SentinelAddrs) > 0 {
		rc.Client, err = (radix.SentinelConfig{
			PoolConfig: radix.PoolConfig{
				Dialer: radix.Dialer{
					SelectDB:	rc.Config.DB,
					AuthUser:	rc.Config.AuthUser,
					AuthPass:	rc.Config.AuthPassword,
					NetDialer:	dialer,
				},
			},
		}).New(ctx, rc.Config.SentinelPrimaryName, rc.Config.SentinelAddrs)
	} else if len(rc.Config.ClusterAddrs) > 0 {
		rc.Client, err = (radix.ClusterConfig{
			PoolConfig: radix.PoolConfig{
				Dialer: radix.Dialer{
					SelectDB:	rc.Config.DB,
					AuthUser:	rc.Config.AuthUser,
					AuthPass:	rc.Config.AuthPassword,
					NetDialer:	dialer,
				},
			},
		}).New(ctx, rc.Config.ClusterAddrs)
	} else {
		rc.Client, err = (radix.PoolConfig{
			Dialer: radix.Dialer{
				SelectDB:	rc.Config.DB,
				AuthUser:	rc.Config.AuthUser,
				AuthPass:	rc.Config.AuthPassword,
				NetDialer:	dialer,
			},
		}).New(ctx, rc.Config.Network, rc.Config.Addr)
	}
	if err != nil {
		return nil, ex.New(err)
	}
	return &rc, nil
}

// Assert `RadixClient` implements `Client`.
var (
	_ Client = (*RadixClient)(nil)
)

// RadixNetDialer is a dialer for radix connections.
type RadixNetDialer interface {
	DialContext(context.Context, string, string) (net.Conn, error)
}

// RadixDoCloser is an thin implementation of the radix client.
type RadixDoCloser interface {
	Do(context.Context, radix.Action) error
	Close() error
}

// RadixClient is a wrapping client for the underling radix redis driver.
type RadixClient struct {
	Config	Config
	Log	logger.Triggerable
	Tracer	Tracer
	Client	RadixDoCloser
}

// Ping sends an echo to the server and validates the response.
func (rc *RadixClient) Ping(ctx context.Context) error {
	var actual string
	expected := uuid.V4().String()
	if err := rc.Client.Do(ctx, radix.Cmd(&actual, OpECHO, expected)); err != nil {
		return ex.New(err)
	}
	if actual != expected {
		return ex.New(ErrPingFailed)
	}
	return nil
}

// Check implements a status check.
func (rc *RadixClient) Check(ctx context.Context) error {
	return rc.Ping(ctx)
}

// Do runs a given command.
func (rc *RadixClient) Do(ctx context.Context, out interface{}, op string, args ...string) (err error) {
	if rc.Log != nil {
		started := time.Now()
		defer func() {
			rc.Log.TriggerContext(ctx, NewEvent(op, args, time.Since(started),
				OptEventNetwork(rc.Config.Network),
				OptEventAddr(rc.Config.Addr),
				OptEventAuthUser(rc.Config.AuthUser),
				OptEventDB(rc.Config.DB),
				OptEventErr(err),
			))
		}()
	}
	if rc.Tracer != nil {
		finisher := rc.Tracer.Do(ctx, rc.Config, op, args)
		defer finisher.Finish(ctx, err)
	}
	if rc.Config.Timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, rc.Config.Timeout)
		defer cancel()
	}
	if radixErr := rc.Client.Do(ctx, radix.Cmd(out, op, args...)); radixErr != nil {
		err = ex.New(radixErr)
		return
	}
	return
}

// Pipeline runs the given commands in a pipeline
func (rc *RadixClient) Pipeline(ctx context.Context, pipelineName string, ops ...Operation) (err error) {
	// Create a parent span for the entire pipeline operation
	if rc.Tracer != nil && pipelineName != "" {
		pipelineFinisher := rc.Tracer.Do(ctx, rc.Config, pipelineName, nil)
		defer pipelineFinisher.Finish(ctx, err)
	}

	var logEvents []Event
	if rc.Log != nil {
		started := time.Now()
		for _, op := range ops {
			event := NewEvent(op.Command, op.Args, time.Since(started),
				OptEventNetwork(rc.Config.Network),
				OptEventAddr(rc.Config.Addr),
				OptEventAuthUser(rc.Config.AuthUser),
				OptEventDB(rc.Config.DB),
				OptEventErr(err),
			)
			logEvents = append(logEvents, event)
		}

		defer func() {
			for _, event := range logEvents {
				rc.Log.TriggerContext(ctx, event)
			}
		}()
	}

	// create a child span for each op
	var finishers []TraceFinisher
	if rc.Tracer != nil {
		for _, op := range ops {
			finisher := rc.Tracer.Do(ctx, rc.Config, op.Command, op.Args)
			finishers = append(finishers, finisher)
		}

		defer func() {
			for _, finisher := range finishers {
				finisher.Finish(ctx, err)
			}
		}()
	}

	if rc.Config.Timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, rc.Config.Timeout)
		defer cancel()
	}

	p := radix.NewPipeline()
	for _, op := range ops {
		p.Append(radix.Cmd(op.Out, op.Command, op.Args...))
	}

	// Execute pipeline
	radixErr := rc.Client.Do(ctx, p)
	if radixErr != nil {
		err = ex.New(radixErr)
		return
	}
	return
}

// Close closes the underlying connection.
func (rc *RadixClient) Close() error {
	return rc.Client.Close()
}
