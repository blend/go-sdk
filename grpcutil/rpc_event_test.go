/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestRPCEvent(t *testing.T) {
	assert := assert.New(t)

	re := NewRPCEvent("/v1.foo", time.Second,
		OptRPCAuthority("event-authority"),
		OptRPCContentType("event-content-type"),
		OptRPCElapsed(time.Millisecond),
		OptRPCEngine("event-engine"),
		OptRPCErr(fmt.Errorf("test error")),
		OptRPCMethod("/v1.bar"),
		OptRPCPeer("event-peer"),
		OptRPCUserAgent("event-user-agent"),
	)

	assert.Equal("event-authority", re.Authority)
	assert.Equal("event-content-type", re.ContentType)
	assert.Equal(time.Millisecond, re.Elapsed)
	assert.Equal("event-engine", re.Engine)
	assert.Equal(fmt.Errorf("test error"), re.Err)
	assert.Equal("/v1.bar", re.Method)
	assert.Equal("event-peer", re.Peer)
	assert.Equal("event-user-agent", re.UserAgent)

	buf := new(bytes.Buffer)
	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}

	re.WriteText(noColor, buf)
	assert.Equal("[event-engine] /v1.bar event-peer event-authority event-user-agent event-content-type 1ms failed", buf.String())

	contents, err := json.Marshal(re)
	assert.Nil(err)
	assert.Contains(string(contents), "event-engine")
}

func TestRPCEvent_StatusCode(t *testing.T) {
	assert := assert.New(t)

	re := NewRPCEvent("/v1.foo", time.Second,
		OptRPCAuthority("event-authority"),
		OptRPCContentType("event-content-type"),
		OptRPCElapsed(time.Millisecond),
		OptRPCEngine("event-engine"),
		OptRPCErr(status.Error(codes.ResourceExhausted, "done goofed kid")),
		OptRPCMethod("/v1.bar"),
		OptRPCPeer("event-peer"),
		OptRPCUserAgent("event-user-agent"),
	)

	buf := new(bytes.Buffer)
	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}

	re.WriteText(noColor, buf)
	assert.Equal("[event-engine] /v1.bar event-peer event-authority event-user-agent event-content-type 1ms failed (8 - ResourceExhausted)", buf.String())

	contents, err := json.Marshal(re)
	assert.Nil(err)
	assert.Contains(string(contents), "event-engine")
}

func TestRPCEventListener(t *testing.T) {
	assert := assert.New(t)

	re := NewRPCEvent("/v1.foo", time.Second)

	var didCall bool
	ml := NewRPCEventListener(func(ctx context.Context, e RPCEvent) {
		didCall = true
	})
	ml(context.Background(), re)
	assert.True(didCall)
}
