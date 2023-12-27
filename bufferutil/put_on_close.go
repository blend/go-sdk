/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package bufferutil

import "bytes"

// PutOnClose wraps a buffer with a close function that will return
// the buffer to the pool.
func PutOnClose(buffer *bytes.Buffer, pool *Pool) *PutOnCloser {
	return &PutOnCloser{Buffer: buffer, Pool: pool}
}

// PutOnCloser is a helper wrapper that will return a buffer to a given pool.
type PutOnCloser struct {
	*bytes.Buffer
	Pool	*Pool
}

// Close returns the buffer to the pool.
func (poc PutOnCloser) Close() error {
	poc.Pool.Put(poc.Buffer)
	return nil
}
