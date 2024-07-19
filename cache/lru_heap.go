/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

import "container/heap"

var (
	_ LRU = (*LRUHeap)(nil)
)

// NewLRUHeap creates a new, empty, LRU Heap.
func NewLRUHeap() *LRUHeap {
	return new(LRUHeap)
}

// LRUHeap is a fifo buffer that is backed by a pre-allocated array, instead of allocating
// a whole new node object for each element (which saves GC churn).
// Enqueue can be O(n), Dequeue can be O(1).
type LRUHeap struct {
	Values LRUHeapValues
}

// Len returns the length of the queue (as it is currently populated).
// Actual memory footprint may be different.
func (lrh *LRUHeap) Len() int {
	return len(lrh.Values)
}

// Push adds an element to the heap.
func (lrh *LRUHeap) Push(object *Value) {
	heap.Push(&lrh.Values, object)
}

// Pop removes the first (oldest) element from the heap.
func (lrh *LRUHeap) Pop() *Value {
	if len(lrh.Values) == 0 {
		return nil
	}
	return heap.Pop(&lrh.Values).(*Value)
}

// Fix updates a value by key.
func (lrh *LRUHeap) Fix(newValue *Value) {
	if len(lrh.Values) == 0 {
		return
	}
	var i int
	for index, value := range lrh.Values {
		if value.Key == newValue.Key {
			i = index
			break
		}
	}
	lrh.Values[i] = newValue
	heap.Fix(&lrh.Values, i)
}

// Remove removes a value by key.
func (lrh *LRUHeap) Remove(key interface{}) {
	if len(lrh.Values) == 0 {
		return
	}
	var i int
	for index, value := range lrh.Values {
		if value.Key == key {
			i = index
			break
		}
	}
	heap.Remove(&lrh.Values, i)
}

// Peek returns the oldest value but does not dequeue it.
func (lrh *LRUHeap) Peek() *Value {
	if len(lrh.Values) == 0 {
		return nil
	}
	return lrh.Values[0]
}

// Consume calls the consumer for each element in the buffer, while also dequeueing that entry.
// The consumer should return `true` if it should remove the item and continue processing.
// If `false` is returned, the current item will be left in place.
func (lrh *LRUHeap) Consume(consumer func(value *Value) bool) {
	if len(lrh.Values) == 0 {
		return
	}

	len := len(lrh.Values)
	for i := 0; i < len; i++ {
		if !consumer(lrh.Peek()) {
			return
		}
		lrh.Pop()
	}
}

// Reset removes all elements from the heap, leaving an empty heap.
func (lrh *LRUHeap) Reset() {
	lrh.Values = nil
}
