/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSyncRingBuffer(t *testing.T) {
	a := assert.New(t)

	buffer := NewSyncRingBuffer[int]()

	buffer.Enqueue(1)
	a.Equal(1, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(1, buffer.PeekBack())

	buffer.Enqueue(2)
	a.Equal(2, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(2, buffer.PeekBack())

	buffer.Enqueue(3)
	a.Equal(3, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(3, buffer.PeekBack())

	buffer.Enqueue(4)
	a.Equal(4, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(4, buffer.PeekBack())

	buffer.Enqueue(5)
	a.Equal(5, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(5, buffer.PeekBack())

	buffer.Enqueue(6)
	a.Equal(6, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(6, buffer.PeekBack())

	buffer.Enqueue(7)
	a.Equal(7, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(7, buffer.PeekBack())

	buffer.Enqueue(8)
	a.Equal(8, buffer.Len())
	a.Equal(1, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value := buffer.Dequeue()
	a.Equal(1, value)
	a.Equal(7, buffer.Len())
	a.Equal(2, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(2, value)
	a.Equal(6, buffer.Len())
	a.Equal(3, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(3, value)
	a.Equal(5, buffer.Len())
	a.Equal(4, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(4, value)
	a.Equal(4, buffer.Len())
	a.Equal(5, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(5, value)
	a.Equal(3, buffer.Len())
	a.Equal(6, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(6, value)
	a.Equal(2, buffer.Len())
	a.Equal(7, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(7, value)
	a.Equal(1, buffer.Len())
	a.Equal(8, buffer.Peek())
	a.Equal(8, buffer.PeekBack())

	value = buffer.Dequeue()
	a.Equal(8, value)
	a.Equal(0, buffer.Len())
	a.Empty(buffer.Peek())
	a.Empty(buffer.PeekBack())
}

func TestSynchronizedRingBufferClear(t *testing.T) {
	a := assert.New(t)

	buffer := NewSyncRingBuffer[int]()
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)

	a.Equal(8, buffer.Len())

	buffer.Clear()
	a.Equal(0, buffer.Len())
	a.Empty(buffer.Peek())
	a.Empty(buffer.PeekBack())
}

func TestSynchronizedRingBufferAsSlice(t *testing.T) {
	a := assert.New(t)

	buffer := NewSyncRingBuffer[int]()
	buffer.Enqueue(1)
	buffer.Enqueue(2)
	buffer.Enqueue(3)
	buffer.Enqueue(4)
	buffer.Enqueue(5)

	contents := buffer.Contents()
	a.Len(contents, 5)
	a.Equal(1, contents[0])
	a.Equal(2, contents[1])
	a.Equal(3, contents[2])
	a.Equal(4, contents[3])
	a.Equal(5, contents[4])
}

func TestSynchronizedRingBufferEach(t *testing.T) {
	a := assert.New(t)

	buffer := NewSyncRingBuffer[int]()

	for x := 1; x < 17; x++ {
		buffer.Enqueue(x)
	}

	var called int
	buffer.Each(func(v int) {
		if v == (called + 1) {
			called++
		}
	})

	a.Equal(16, called)
}

func TestSynchronizedRingBufferDrain(t *testing.T) {
	a := assert.New(t)

	buffer := NewSyncRingBuffer[int]()

	for x := 1; x < 17; x++ {
		buffer.Enqueue(x)
	}

	a.Equal(16, buffer.Len())

	var called int
	buffer.Consume(func(v int) {
		called++
	})

	a.Equal(16, called)
	a.Zero(buffer.Len())
}
