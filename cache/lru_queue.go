/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

import (
	"sort"
)

var (
	_ LRU = (*LRUQueue)(nil)
)

// NewLRUQueue creates a new, empty, LRUQueue.
func NewLRUQueue() *LRUQueue {
	return &LRUQueue{
		array: make([]*Value, ringBufferDefaultCapacity),
	}
}

// LRUQueue is a fifo buffer that is backed by a pre-allocated array, instead of allocating
// a whole new node object for each element (which saves GC churn).
// Enqueue can be O(n), Dequeue can be O(1).
type LRUQueue struct {
	array []*Value
	head  int
	tail  int
	size  int
}

// Len returns the length of the ring buffer (as it is currently populated).
// Actual memory footprint may be different.
func (lru *LRUQueue) Len() (len int) {
	return lru.size
}

// Capacity returns the total size of the ring bufffer, including empty elements.
func (lru *LRUQueue) Capacity() int {
	return len(lru.array)
}

// Clear removes all objects from the LRUQueue.
func (lru *LRUQueue) Clear() {
	if lru.head < lru.tail {
		arrayClear(lru.array, lru.head, lru.size)
	} else {
		arrayClear(lru.array, lru.head, len(lru.array)-lru.head)
		arrayClear(lru.array, 0, lru.tail)
	}
	lru.head = 0
	lru.tail = 0
	lru.size = 0
}

// Push adds an element to the "back" of the LRUQueue.
func (lru *LRUQueue) Push(object *Value) {
	if lru.size == len(lru.array) { // if we're out of room
		lru.setCapacity(lru.growCapacity())
	}
	lru.array[lru.tail] = object
	lru.tail = (lru.tail + 1) % len(lru.array)
	lru.size++
}

// Pop removes the first (oldest) element from the LRUQueue.
func (lru *LRUQueue) Pop() *Value {
	if lru.size == 0 {
		return nil
	}

	removed := lru.array[lru.head]
	lru.head = (lru.head + 1) % len(lru.array)
	lru.size--
	return removed
}

// Peek returns but does not remove the first element.
func (lru *LRUQueue) Peek() *Value {
	if lru.size == 0 {
		return nil
	}
	return lru.array[lru.head]
}

// PeekBack returns but does not remove the last element.
func (lru *LRUQueue) PeekBack() *Value {
	if lru.size == 0 {
		return nil
	}
	if lru.tail == 0 {
		return lru.array[len(lru.array)-1]
	}
	return lru.array[lru.tail-1]
}

// Fix updates the queue given an update to a specific value.
func (lru *LRUQueue) Fix(value *Value) {
	if lru.size == 0 {
		return
	}
	if value == nil {
		panic("lru queue; value is nil")
	}

	values := make([]*Value, lru.size)
	var index int
	var didUpdate bool
	lru.Each(func(v *Value) bool {
		if v.Key == value.Key {
			didUpdate = v.Expires != value.Expires
			values[index] = value
		} else {
			values[index] = v
		}
		index++
		return true
	})
	if didUpdate {
		sort.Sort(LRUHeapValues(values))
	}
	lru.array = make([]*Value, len(lru.array))
	copy(lru.array, values)
	lru.head = 0
	lru.tail = lru.size
}

// Remove removes an item from the queue by its key.
func (lru *LRUQueue) Remove(key interface{}) {
	if lru.size == 0 {
		return
	}
	if key == nil {
		panic("lru queue; key is nil")
	}

	size := lru.size

	values := make([]*Value, size-1)
	var cursor int
	for x := 0; x < size; x++ {
		head := lru.Pop()
		if head.Key != key {
			values[cursor] = head
			cursor++
		}
	}
	for x := 0; x < len(values); x++ {
		lru.Push(values[x])
	}
}

// Each iterates through the queue and calls the consumer for each element of the queue.
func (lru *LRUQueue) Each(consumer func(*Value) bool) {
	if lru.size == 0 {
		return
	}

	if lru.head < lru.tail {
		for cursor := lru.head; cursor < lru.tail; cursor++ {
			if !consumer(lru.array[cursor]) {
				return
			}
		}
	} else {
		for cursor := lru.head; cursor < len(lru.array); cursor++ {
			if !consumer(lru.array[cursor]) {
				return
			}
		}
		for cursor := 0; cursor < lru.tail; cursor++ {
			if !consumer(lru.array[cursor]) {
				return
			}
		}
	}
}

// Consume calls the consumer for each element in the buffer. If the handler returns true,
// the element is popped and the handler is called on the next value.
func (lru *LRUQueue) Consume(consumer func(*Value) bool) {
	if lru.size == 0 {
		return
	}

	for i := 0; i < lru.size; i++ {
		if !consumer(lru.Peek()) {
			return
		}
		lru.Pop()
	}
}

// Reset removes all elements from the heap, leaving an empty heap.
func (lru *LRUQueue) Reset() {
	lru.array = make([]*Value, ringBufferDefaultCapacity)
	lru.head = 0
	lru.tail = 0
	lru.size = 0
}

//
// util / helpers
//

// TrimExcess trims the excess space in the ringbuffer.
func (lru *LRUQueue) TrimExcess() {
	threshold := float64(len(lru.array)) * 0.9
	if lru.size < int(threshold) {
		lru.setCapacity(lru.size)
	}
}

//
// internal helpers
//

func (lru *LRUQueue) growCapacity() int {
	size := len(lru.array)
	newCapacity := size << 1
	minimumGrow := size + ringBufferMinimumGrow

	if newCapacity < minimumGrow {
		newCapacity = minimumGrow
	}

	return newCapacity
}

func (lru *LRUQueue) setCapacity(capacity int) {
	newArray := make([]*Value, capacity)
	if lru.size > 0 {
		if lru.head < lru.tail {
			arrayCopy(lru.array, lru.head, newArray, 0, lru.size)
		} else {
			arrayCopy(lru.array, lru.head, newArray, 0, len(lru.array)-lru.head)
			arrayCopy(lru.array, 0, newArray, len(lru.array)-lru.head, lru.tail)
		}
	}
	lru.array = newArray
	lru.head = 0
	if lru.size == capacity {
		lru.tail = 0
	} else {
		lru.tail = lru.size
	}
}

//
// array helpers
//

func arrayClear(source []*Value, index, length int) {
	for x := 0; x < length; x++ {
		absoluteIndex := x + index
		source[absoluteIndex] = nil
	}
}

func arrayCopy(source []*Value, sourceIndex int, destination []*Value, destinationIndex, length int) {
	for x := 0; x < length; x++ {
		from := sourceIndex + x
		to := destinationIndex + x

		destination[to] = source[from]
	}
}

const (
	ringBufferMinimumGrow     = 4
	ringBufferDefaultCapacity = 4
)
