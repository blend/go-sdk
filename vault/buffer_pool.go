/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"bytes"
	"sync"
)

// NewBufferPool returns a new BufferPool.
// bufferSize is the size of the returned buffers pre-allocated size in bytes.
// Typically this is something between 256 bytes and 1kb.
func NewBufferPool(bufferSize int) *BufferPool {
	bp := &BufferPool{}
	bp.Pool = sync.Pool{
		New: func() interface{} {
			b := &Buffer{
				Buffer:	bytes.NewBuffer(make([]byte, bufferSize)),
				pool:	bp,
			}
			return b
		},
	}
	return bp
}

// BufferPool is a sync.Pool of bytes.Buffer.
type BufferPool struct {
	sync.Pool
}

// Get returns a pooled bytes.Buffer instance.
func (bp *BufferPool) Get() *Buffer {
	buf := bp.Pool.Get().(*Buffer)
	buf.Reset()
	return buf
}

// Put returns the pooled instance.
func (bp *BufferPool) Put(b *Buffer) {
	bp.Pool.Put(b)
}

// Buffer is a bytes.Buffer with a reference back to the buffer pool.
// It returns itself to the pool on close.
type Buffer struct {
	*bytes.Buffer
	pool	*BufferPool
}

// Close returns the buffer to the pool.
func (b *Buffer) Close() error {
	b.Reset()
	b.pool.Put(b)
	return nil
}
