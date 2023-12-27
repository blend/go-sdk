/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	radix "github.com/mediocregopher/radix/v4"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/redis"
)

func Test_RadixClient_Do(t *testing.T) {
	its := assert.New(t)

	buf := new(bytes.Buffer)
	log := logger.Memory(buf)
	defer log.Close()

	logEvents := make(chan redis.Event)
	log.Listen("test", "test", redis.NewEventListener(func(_ context.Context, e redis.Event) {
		logEvents <- e
	}))

	mockRadixClient := &MockRadixClient{
		Ops: make(chan radix.Action, 1),
	}

	rc := &redis.RadixClient{
		Log:	log,
		Client:	mockRadixClient,
	}

	var foo string
	its.Nil(rc.Do(context.TODO(), &foo, "GET", "foo"))
}

func Test_RadixClient_Do_timeout(t *testing.T) {
	its := assert.New(t)

	mockRadixClient := &MockRadixClient{
		Ops: make(chan radix.Action),
	}
	rc := &redis.RadixClient{
		Config: redis.Config{
			Timeout: time.Millisecond,
		},
		Client:	mockRadixClient,
	}
	var foo string
	its.NotNil(rc.Do(context.Background(), &foo, "GET", "foo"))
}

func Test_RadixClient_Pipeline(t *testing.T) {
	its := assert.New(t)

	// Mock tracer and finisher
	mockTracerFinisher := &MockTraceFinisher{}
	mockTracer := &MockTracer{mockTraceFinisher: mockTracerFinisher}
	log := &MockTriggerable{}

	// Mock RadixClient
	mockRadixClient := &MockRadixClient{
		Ops: make(chan radix.Action, 10),
	}

	rc := &redis.RadixClient{
		Log:	log,
		Tracer:	mockTracer,
		Client:	mockRadixClient,
	}

	// Mock operations
	ops := []redis.Operation{
		{Command: "GET", Args: []string{"key1"}},
		{Command: "SET", Args: []string{"key2", "value"}},
	}

	// Perform Pipeline operation
	err := rc.Pipeline(context.TODO(), "GetSetPipeline", ops...)

	// Assertions
	its.Nil(err)
	// 1 (for parent span) + 2 (for child spans)
	its.Equal(3, len(mockTracer.calls))
	its.Equal(3, len(mockTracerFinisher.calls))
}

// MockTracer is a mock of Tracer that stores calls in-memory
type MockTracer struct {
	calls	[]struct {
		Op	string
		Args	[]string
	}
	mockTraceFinisher	*MockTraceFinisher
}

// Do implements Tracer
func (mt *MockTracer) Do(ctx context.Context, cfg redis.Config, op string, args []string) redis.TraceFinisher {
	mt.calls = append(mt.calls, struct {
		Op	string
		Args	[]string
	}{op, args})
	return mt.mockTraceFinisher
}

// MockTracerFinisher is a mock of TraceFinisher
type MockTraceFinisher struct {
	calls []struct {
		Err error
	}
}

// Finish implements Tracer
func (mtf *MockTraceFinisher) Finish(ctx context.Context, err error) {
	mtf.calls = append(mtf.calls, struct {
		Err error
	}{err})
}

// MockTriggerable is a of a mock of Triggerable that stores calls in-memory
type MockTriggerable struct {
	Events []struct {
		Event logger.Event
	}
}

// TriggerContext implements Triggerable
func (mt *MockTriggerable) TriggerContext(ctx context.Context, e logger.Event) {
	mt.Events = append(mt.Events, struct {
		Event logger.Event
	}{e})
}
