/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package shardutil

import (
	"context"
	"sync"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

// Shards handles communicating with many underlying databases at once.
type Shards struct {
	Connections []*db.Connection
	Opts        []InvocationOption
}

// PartitionIndex returns a partition index for a given hashCode.
func (s Shards) PartitionIndex(hashCode int) int {
	return hashCode % len(s.Connections)
}

// PartitionOptions returns db.InvocationOptions for a given partition.
func (s Shards) PartitionOptions(partitionIndex int, opts ...InvocationOption) []db.InvocationOption {
	var invocationOpts []db.InvocationOption
	for _, opt := range s.Opts {
		invocationOpts = append(invocationOpts, opt(partitionIndex))
	}
	for _, opt := range opts {
		invocationOpts = append(invocationOpts, opt(partitionIndex))
	}
	return invocationOpts
}

// InvokeAll invokes a given function asynchronously for each connection in the manager.
func (s Shards) InvokeAll(ctx context.Context, action func(int, *db.Invocation) error, opts ...InvocationOption) error {
	wg := new(sync.WaitGroup)
	wg.Add(len(s.Connections))

	errors := make(chan error, len(s.Connections))
	for index := 0; index < len(s.Connections); index++ {
		go func(partitionIndex int) {
			defer func() {
				if r := recover(); r != nil {
					errors <- ex.New(r)
				}
				wg.Done()
			}()

			invocation := s.Connections[partitionIndex].Invoke(
				append(s.PartitionOptions(partitionIndex, opts...), db.OptContext(ctx))...,
			)
			if err := action(partitionIndex, invocation); err != nil {
				errors <- err
			}
		}(index)
	}

	wg.Wait()
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}

// InvokeOne creates a new db invocation routed to an underlying connection mapped by a given hashcode.
// The underlying connection is determined by `PartitionIndex(hashCode)`.
// The options are special parameterized versions of normal `db.InvocationOptions` that also take a partition index.
// The returned invocation will map to only (1) underlying connection.
func (s Shards) InvokeOne(ctx context.Context, hashCode int, opts ...InvocationOption) *db.Invocation {
	partitionIndex := s.PartitionIndex(hashCode)
	return s.Connections[partitionIndex].Invoke(append(s.PartitionOptions(partitionIndex, opts...), db.OptContext(ctx))...)
}
