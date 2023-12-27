/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import (
	"context"
	"io"
)

// Client is the basic interface that a redis client should implement.
type Client interface {
	io.Closer
	Do(ctx context.Context, out interface{}, command string, args ...string) error
	Pipeline(ctx context.Context, pipelineName string, ops ...Operation) error
}

// Operation encapsulates a redis command to be made to the client
type Operation struct {
	Out	interface{}
	Command	string
	Args	[]string
}
