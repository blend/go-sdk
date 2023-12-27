/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

import (
	"sync"
)

// NewSyncRingBuffer returns a new synchronized ring buffer.
func NewSyncRingBuffer[T any]() *SyncRingBuffer[T] {
	return &SyncRingBuffer[T]{
		innerBuffer:	NewRingBuffer[T](),
		syncRoot:	&sync.Mutex{},
	}
}

// NewSyncRingBufferWithCapacity returns a new synchronized ring buffer.
func NewSyncRingBufferWithCapacity[T any](capacity int) *SyncRingBuffer[T] {
	return &SyncRingBuffer[T]{
		innerBuffer:	NewRingBufferWithCapacity[T](capacity),
		syncRoot:	&sync.Mutex{},
	}
}

// SyncRingBuffer is a ring buffer wrapper that adds synchronization.
type SyncRingBuffer[T any] struct {
	innerBuffer	*RingBuffer[T]
	syncRoot	*sync.Mutex
}

// SyncRoot returns the mutex used to synchronize the collection.
func (srb *SyncRingBuffer[T]) SyncRoot() *sync.Mutex {
	return srb.syncRoot
}

// RingBuffer returns the inner ring buffer.
func (srb *SyncRingBuffer[T]) RingBuffer() *RingBuffer[T] {
	return srb.innerBuffer
}

// Len returns the length of the ring buffer (as it is currently populated).
// Actual memory footprint may be different.
func (srb SyncRingBuffer[T]) Len() int {
	srb.syncRoot.Lock()
	defer srb.syncRoot.Unlock()
	return srb.innerBuffer.Len()
}

// Capacity returns the total size of the ring buffer, including empty elements.
func (srb *SyncRingBuffer[T]) Capacity() int {
	srb.syncRoot.Lock()
	defer srb.syncRoot.Unlock()
	return srb.innerBuffer.Capacity()
}

// Enqueue adds an element to the "back" of the ring buffer.
func (srb *SyncRingBuffer[T]) Enqueue(value T) {
	srb.syncRoot.Lock()
	srb.innerBuffer.Enqueue(value)
	srb.syncRoot.Unlock()
}

// Dequeue removes the first (oldest) element from the ring buffer.
func (srb *SyncRingBuffer[T]) Dequeue() T {
	var val T
	srb.syncRoot.Lock()
	val = srb.innerBuffer.Dequeue()
	srb.syncRoot.Unlock()
	return val
}

// DequeueBack removes the last (newest) element from the ring buffer.
func (srb *SyncRingBuffer[T]) DequeueBack() T {
	var val T
	srb.syncRoot.Lock()
	val = srb.innerBuffer.DequeueBack()
	srb.syncRoot.Unlock()
	return val
}

// Peek returns but does not remove the first element.
func (srb *SyncRingBuffer[T]) Peek() T {
	var val T
	srb.syncRoot.Lock()
	val = srb.innerBuffer.Peek()
	srb.syncRoot.Unlock()
	return val
}

// PeekBack returns but does not remove the last element.
func (srb *SyncRingBuffer[T]) PeekBack() T {
	var val T
	srb.syncRoot.Lock()
	val = srb.innerBuffer.PeekBack()
	srb.syncRoot.Unlock()
	return val
}

// TrimExcess resizes the buffer to better fit the contents.
func (srb *SyncRingBuffer[T]) TrimExcess() {
	srb.syncRoot.Lock()
	srb.innerBuffer.trimExcess()
	srb.syncRoot.Unlock()
}

// Contents returns the ring buffer, in order, as a slice.
func (srb *SyncRingBuffer[T]) Contents() []T {
	var val []T
	srb.syncRoot.Lock()
	val = srb.innerBuffer.Contents()
	srb.syncRoot.Unlock()
	return val
}

// Clear removes all objects from the ring buffer.
func (srb *SyncRingBuffer[T]) Clear() {
	srb.syncRoot.Lock()
	srb.innerBuffer.Clear()
	srb.syncRoot.Unlock()
}

// Drain returns the ring buffer, in order, as a slice and empties it.
func (srb *SyncRingBuffer[T]) Drain() []T {
	var val []T
	srb.syncRoot.Lock()
	val = srb.innerBuffer.Drain()
	srb.syncRoot.Unlock()
	return val
}

// Each calls the consumer for each element in the buffer.
func (srb *SyncRingBuffer[T]) Each(consumer func(value T)) {
	srb.syncRoot.Lock()
	srb.innerBuffer.Each(consumer)
	srb.syncRoot.Unlock()
}

// Consume calls the consumer for each element in the buffer, while also dequeueing that entry.
func (srb *SyncRingBuffer[T]) Consume(consumer func(value T)) {
	srb.syncRoot.Lock()
	srb.innerBuffer.Consume(consumer)
	srb.syncRoot.Unlock()
}

// EachUntil calls the consumer for each element in the buffer with a stopping condition in head=>tail order.
func (srb *SyncRingBuffer[T]) EachUntil(consumer func(value T) bool) {
	srb.syncRoot.Lock()
	srb.innerBuffer.EachUntil(consumer)
	srb.syncRoot.Unlock()
}

// ReverseEachUntil calls the consumer for each element in the buffer with a stopping condition in tail=>head order.
func (srb *SyncRingBuffer[T]) ReverseEachUntil(consumer func(value T) bool) {
	srb.syncRoot.Lock()
	srb.innerBuffer.ReverseEachUntil(consumer)
	srb.syncRoot.Unlock()
}
