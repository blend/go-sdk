/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"context"
	"time"

	"google.golang.org/grpc/connectivity"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
)

// CheckConnectivityState returns an async checker for a client that provides a connection state.
func CheckConnectivityState(client ConnectionStateProvider, opts ...RetryCheckConnectivityStateOption) async.Checker {
	return async.CheckerFunc(func(ctx context.Context) error {
		state, err := RetryCheckConnectivityState(ctx, client, opts...)
		if err != nil {
			return err
		}
		if state != connectivity.Ready {
			return ex.New(ErrConnectionNotReady)
		}
		return nil
	})
}

// ConnectionStateProvider is a type that provides a connection state.
type ConnectionStateProvider interface {
	GetConnectionState() connectivity.State
}

// RetryCheckConnectivityStateOptions are options for checking the connectivity state.
type RetryCheckConnectivityStateOptions struct {
	RetryTimeout time.Duration
	RetryBackoff time.Duration
	MaxRetries   uint
}

// RetryCheckConnectivityStateOption mutates CheckConnectivityStateOptions.
type RetryCheckConnectivityStateOption func(*RetryCheckConnectivityStateOptions)

// OptRetryCheckConnectivityStateRetryTimeout sets the RetryTimeout.
func OptRetryCheckConnectivityStateRetryTimeout(d time.Duration) RetryCheckConnectivityStateOption {
	return func(opts *RetryCheckConnectivityStateOptions) {
		opts.RetryTimeout = d
	}
}

// OptRetryCheckConnectivityStateRetryBackoff sets the RetryBackoff.
func OptRetryCheckConnectivityStateRetryBackoff(d time.Duration) RetryCheckConnectivityStateOption {
	return func(opts *RetryCheckConnectivityStateOptions) {
		opts.RetryBackoff = d
	}
}

// OptRetryCheckConnectivityStateMaxRetries sets the MaxRetries.
func OptRetryCheckConnectivityStateMaxRetries(maxRetries uint) RetryCheckConnectivityStateOption {
	return func(opts *RetryCheckConnectivityStateOptions) {
		opts.MaxRetries = maxRetries
	}
}

// ErrConnectionNotReady is returned by ConnectivityStateChecker.
const ErrConnectionNotReady ex.Class = "grpc connection not ready"

// RetryCheckConnectivityState implements a retry checker for connectivity state.
func RetryCheckConnectivityState(ctx context.Context, client ConnectionStateProvider, opts ...RetryCheckConnectivityStateOption) (state connectivity.State, err error) {
	options := RetryCheckConnectivityStateOptions{
		RetryTimeout: 5 * time.Second,
		RetryBackoff: 200 * time.Millisecond,
		MaxRetries:   30,
	}
	for _, opt := range opts {
		opt(&options)
	}

	state = client.GetConnectionState()
	if state != connectivity.Ready {
		alarm := time.NewTimer(options.RetryTimeout)
		tick := time.NewTicker(options.RetryBackoff)
		defer alarm.Stop()
		defer tick.Stop()

		for state != connectivity.Ready {
			select {
			case <-tick.C:
				state = client.GetConnectionState()
			case <-alarm.C:
				state = client.GetConnectionState()
				return
			case <-ctx.Done():
				state = client.GetConnectionState()
				err = context.Canceled
				return
			}
		}
	}
	return
}
