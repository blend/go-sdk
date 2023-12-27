/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestLRUHeap(t *testing.T) {
	assert := assert.New(t)

	t0 := time.Date(2019, 06, 13, 12, 10, 9, 8, time.UTC)
	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 15, 12, 10, 9, 8, time.UTC)
	t3 := time.Date(2019, 06, 16, 12, 10, 9, 8, time.UTC)
	t4 := time.Date(2019, 06, 17, 12, 10, 9, 8, time.UTC)
	t5 := time.Date(2019, 06, 18, 12, 10, 9, 8, time.UTC)

	h := NewLRUHeap()
	h.Push(&Value{
		Key:		"5",
		Expires:	t5,
	})
	assert.Len(h.Values, 1)
	assert.Equal(1, h.Len())
	h.Push(&Value{
		Key:		"2",
		Expires:	t2,
	})
	assert.Len(h.Values, 2)
	h.Push(&Value{
		Key:		"3",
		Expires:	t3,
	})
	assert.Len(h.Values, 3)
	h.Push(&Value{
		Key:		"0",
		Expires:	t0,
	})
	assert.Len(h.Values, 4)
	h.Push(&Value{
		Key:		"4",
		Expires:	t4,
	})
	assert.Len(h.Values, 5)
	h.Push(&Value{
		Key:		"1",
		Expires:	t1,
	})
	assert.Len(h.Values, 6)
	assert.Equal(t0, h.Values[0].Expires)

	popped := h.Pop()
	assert.Equal(t0, popped.Expires)
	assert.Equal(t1, h.Values[0].Expires)

	popped = h.Pop()
	assert.Equal(t1, popped.Expires)
	assert.Equal(t2, h.Values[0].Expires)

	popped = h.Pop()
	assert.Equal(t2, popped.Expires)
	assert.Equal(t3, h.Values[0].Expires)

	popped = h.Pop()
	assert.Equal(t3, popped.Expires)
	assert.Equal(t4, h.Values[0].Expires)

	popped = h.Pop()
	assert.Equal(t4, popped.Expires)
	assert.Equal(t5, h.Values[0].Expires)

	popped = h.Pop()
	assert.Equal(t5, popped.Expires)
	assert.Empty(h.Values)

	popped = h.Pop()
	assert.Nil(popped)
}

func TestLRUHeapEmpty(t *testing.T) {
	assert := assert.New(t)

	h := NewLRUHeap()
	h.Fix(nil)
	h.Remove(nil)
	assert.Nil(h.Pop())
	assert.Nil(h.Peek())
}

func TestLRUHeapConsumeUntil(t *testing.T) {
	assert := assert.New(t)

	t0 := time.Date(2019, 06, 13, 12, 10, 9, 8, time.UTC)
	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 15, 12, 10, 9, 8, time.UTC)
	t3 := time.Date(2019, 06, 16, 12, 10, 9, 8, time.UTC)
	t4 := time.Date(2019, 06, 17, 12, 10, 9, 8, time.UTC)
	t5 := time.Date(2019, 06, 18, 12, 10, 9, 8, time.UTC)

	h := NewLRUHeap()
	h.Push(&Value{Key: "5", Expires: t5})
	h.Push(&Value{Key: "2", Expires: t2})
	h.Push(&Value{Key: "3", Expires: t3})
	h.Push(&Value{Key: "0", Expires: t0})
	h.Push(&Value{Key: "4", Expires: t4})
	h.Push(&Value{Key: "1", Expires: t1})
	assert.Len(h.Values, 6)

	h.Consume(func(v *Value) bool {
		return v.Expires.Before(t3)
	})
	assert.Len(h.Values, 3, "consumeUntil should have removed (3) items")
}

func TestLRUHeapFix(t *testing.T) {
	assert := assert.New(t)

	t0 := time.Date(2019, 06, 13, 12, 10, 9, 8, time.UTC)
	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 15, 12, 10, 9, 8, time.UTC)
	t3 := time.Date(2018, 06, 15, 12, 10, 9, 8, time.UTC)

	h := NewLRUHeap()
	h.Push(&Value{Key: "1", Expires: t0})
	h.Push(&Value{Key: "2", Expires: t1})
	h.Push(&Value{Key: "3", Expires: t2})
	assert.Equal(t0, h.Peek().Expires)

	// do ths fix
	h.Fix(&Value{Key: "3", Expires: t3})

	assert.Equal(t3, h.Peek().Expires)
}

func TestLRUHeapRemove(t *testing.T) {
	assert := assert.New(t)

	t0 := time.Date(2019, 06, 13, 12, 10, 9, 8, time.UTC)
	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 15, 12, 10, 9, 8, time.UTC)

	h := NewLRUHeap()
	h.Push(&Value{Key: "1", Expires: t0})
	h.Push(&Value{Key: "2", Expires: t1})
	h.Push(&Value{Key: "3", Expires: t2})
	assert.Equal(t0, h.Peek().Expires)

	h.Remove("1")
	assert.Equal(t1, h.Peek().Expires)
}
