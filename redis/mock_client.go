/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import "context"

// NewMockClient returns a new mock client with a given capacity.
func NewMockClient(capacity int) *MockClient {
	return &MockClient{
		Ops: make(chan MockClientOp, capacity),
	}
}

// Assert `MockClient` implements client.
var (
	_ Client = (*MockClient)(nil)
)

// MockClient is a mocked client.
type MockClient struct {
	DoMock	func(context.Context, interface{}, string, ...string) error
	Ops	chan MockClientOp
}

// Do applies a command.
func (mc *MockClient) Do(ctx context.Context, out interface{}, op string, args ...string) error {
	mc.Ops <- MockClientOp{Out: out, Op: op, Args: args}
	if mc.DoMock != nil {
		return mc.DoMock(ctx, out, op, args...)
	}
	return nil
}

// Pipeline applies commands in a piipeline.
func (mc *MockClient) Pipeline(_ context.Context, pipelineName string, ops ...Operation) error {
	for _, op := range ops {
		mc.Ops <- MockClientOp{Out: op.Out, Op: op.Command, Args: op.Args}
	}
	return nil
}

// Close closes the mock client.
func (mc *MockClient) Close() error	{ return nil }

// MockClientOp is a mocked client op.
type MockClientOp struct {
	Out	interface{}
	Op	string
	Args	[]string
}
