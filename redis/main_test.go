/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis_test

import (
	"context"

	radix "github.com/mediocregopher/radix/v4"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/redis"
)

// MockRadixClient implements radix.Client for testing.
type MockRadixClient struct {
	radix.Client
	Ops         chan radix.Action
	PipelineOps [][]redis.Operation

	Log    logger.Triggerable
	Tracer redis.Tracer
}

// Do implements part of the radix client interface.
func (mrc *MockRadixClient) Do(ctx context.Context, action radix.Action) error {
	pushDone := make(chan struct{})
	go func() {
		defer close(pushDone)
		mrc.Ops <- action
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-pushDone:
		return nil
	}
}

// Pipeline mimics the Pipeline method for the mock.
func (mrc *MockRadixClient) Pipeline(ctx context.Context, ops ...redis.Operation) error {
	mrc.PipelineOps = append(mrc.PipelineOps, ops)
	return nil
}
