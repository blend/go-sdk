/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cache

// LRU is a type that implements the LRU methods.
type LRU interface {
	// Len returns the number of items in the queue.
	Len() int
	// Push should add a new value. The new minimum value should be returned by `Peek()` and `Pop()`.
	Push(*Value)
	// Pop should remove and return the minimum value, reordering the heap
	// to set a new minimum value.
	Pop() *Value
	// Peek should return (but not remove) the minimum value.
	Peek() *Value
	// Fix should update the LRU, replacing any existing values, reordering the heap.
	Fix(*Value)
	// Remove should remove a value with a given key, compacting the heap.
	Remove(interface{})
	// Consume should iterate through the values. If `true` is removed by the handler,
	// the current value will be removed and the handler will be called on the next value.
	Consume(func(*Value) bool)
	// Reset should remove all values from the LRU, leaving an empty LRU.
	Reset()
}
