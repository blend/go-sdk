/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package bufferutil

import (
	"bytes"
	"sync"
)

/*
NewPool returns a new Pool, which returns bytes buffers pre-sized to a given minimum size.

The purpose of a buffer pool is to reduce the number of gc collections incurred when using bytes buffers
repeatedly; instead of marking the buffer as to be collected, it is returned to the pool to be re-used.

Example:

	pool := bufferutil.NewPool(1024) // pre-allocate 1024 bytes per buffer.

	func() {
		buf := pool.Get()
		defer pool.Put(buf)

		// do things with the buffer ...
	}()
*/
func NewPool(bufferSize int) *Pool {
	return &Pool{
		Pool: sync.Pool{New: func() interface{} {
			b := bytes.NewBuffer(make([]byte, bufferSize))
			b.Reset()
			return b
		}},
	}
}

// Pool is a sync.Pool of bytes.Buffer.
type Pool struct {
	sync.Pool
}

// Get returns a pooled bytes.Buffer instance.
func (p *Pool) Get() *bytes.Buffer {
	return p.Pool.Get().(*bytes.Buffer)
}

// Put returns the pooled instance.
func (p *Pool) Put(b *bytes.Buffer) {
	b.Reset()
	p.Pool.Put(b)
}
