/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func tv(index int) *Value {
	return &Value{
		Key:       index,
		Value:     strconv.Itoa(index),
		Timestamp: time.Now().UTC(),
		Expires:   time.Date(2019, 06, 14, 12, index, 0, 0, time.UTC),
	}
}

func TestLRUQueue(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()

	buffer.Push(tv(1))
	assert.Equal(1, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(1, buffer.PeekBack().Key)

	buffer.Push(tv(2))
	assert.Equal(2, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(2, buffer.PeekBack().Key)

	buffer.Push(tv(3))
	assert.Equal(3, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(3, buffer.PeekBack().Key)

	buffer.Push(tv(4))
	assert.Equal(4, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(4, buffer.PeekBack().Key)

	buffer.Push(tv(5))
	assert.Equal(5, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(5, buffer.PeekBack().Key)

	buffer.Push(tv(6))
	assert.Equal(6, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(6, buffer.PeekBack().Key)

	buffer.Push(tv(7))
	assert.Equal(7, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(7, buffer.PeekBack().Key)

	buffer.Push(tv(8))
	assert.Equal(8, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value := buffer.Pop()
	assert.Equal(1, value.Key)
	assert.Equal(7, buffer.Len())
	assert.Equal(2, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(2, value.Key)
	assert.Equal(6, buffer.Len())
	assert.Equal(3, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(3, value.Key)
	assert.Equal(5, buffer.Len())
	assert.Equal(4, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(4, value.Key)
	assert.Equal(4, buffer.Len())
	assert.Equal(5, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(5, value.Key)
	assert.Equal(3, buffer.Len())
	assert.Equal(6, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(6, value.Key)
	assert.Equal(2, buffer.Len())
	assert.Equal(7, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(7, value.Key)
	assert.Equal(1, buffer.Len())
	assert.Equal(8, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(8, value.Key)
	assert.Equal(0, buffer.Len())
	assert.Nil(buffer.Peek())
	assert.Nil(buffer.PeekBack())

	before := buffer.Capacity()
	buffer.TrimExcess()
	after := buffer.Capacity()
	assert.True(before > after)

}

func TestLRUQueueClear(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()
	for x := 0; x < 8; x++ {
		buffer.Push(tv(x))
	}
	assert.Equal(8, buffer.Len())
	buffer.Clear()
	assert.Equal(0, buffer.Len())
	assert.Nil(buffer.Peek())
	assert.Nil(buffer.PeekBack())
}

func TestLRUQueueFix(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()
	for x := 0; x < 8; x++ {
		buffer.Push(tv(x))
	}

	tv4e := time.Date(2018, 06, 14, 12, 4, 0, 0, time.UTC)
	tv4 := tv(4)
	tv4.Expires = tv4e

	// fix tv4 to have the earliest expires
	buffer.Fix(tv4)

	// it should now be the head of the queue
	head := buffer.Peek()
	assert.Equal(tv4e, head.Expires, fmt.Sprintf("Head is actually %v", head.Key))

	tv6e := time.Date(2017, 06, 14, 12, 4, 0, 0, time.UTC)
	tv6 := tv(6)
	tv6.Expires = tv6e

	buffer.Fix(tv6)

	head = buffer.Peek()
	assert.Equal(tv6e, head.Expires)
}

func TestLRUQueueRemove(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()
	for x := 0; x < 8; x++ {
		buffer.Push(tv(x))
	}

	buffer.Remove(4)

	var didFind bool
	buffer.Each(func(v *Value) bool {
		if v.Key == 4 {
			didFind = true
			return false
		}
		return true
	})
	assert.False(didFind)

	assert.Equal(tv(0).Expires, buffer.Peek().Expires)
	assert.Equal(tv(7).Expires, buffer.PeekBack().Expires)
}

func TestLRUQueueConsume(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()

	for x := 0; x < 20; x++ {
		buffer.Push(tv(x))
	}

	var called int
	buffer.Consume(func(v *Value) bool {
		// stop after 10
		if v.Key.(int) > 10 {
			return false
		}

		if v.Key.(int) == called {
			called++
		}
		return true
	})

	assert.Equal(10, called)
}
