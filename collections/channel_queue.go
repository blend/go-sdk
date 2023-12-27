/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

import "sync"

// NewChannelQueueWithCapacity returns a new ChannelQueue instance.
func NewChannelQueueWithCapacity[T any](capacity int) *ChannelQueue[T] {
	return &ChannelQueue[T]{Capacity: capacity, storage: make(chan T, capacity), latch: sync.Mutex{}}
}

// NewChannelQueueFromValues returns a new ChannelQueue from a given slice of values.
func NewChannelQueueFromValues[T any](values []T) *ChannelQueue[T] {
	capacity := len(values)
	cq := &ChannelQueue[T]{Capacity: capacity, storage: make(chan T, capacity), latch: sync.Mutex{}}
	for _, v := range values {
		cq.storage <- v
	}
	return cq
}

// ChannelQueue is a thread safe queue.
type ChannelQueue[T any] struct {
	Capacity	int
	storage		chan T
	latch		sync.Mutex
}

// Len returns the number of items in the queue.
func (cq *ChannelQueue[T]) Len() int {
	return len(cq.storage)
}

// Enqueue adds an item to the queue.
func (cq *ChannelQueue[T]) Enqueue(item T) {
	cq.storage <- item
}

// Dequeue returns the next element in the queue.
func (cq *ChannelQueue[T]) Dequeue() T {
	var res T
	if len(cq.storage) != 0 {
		return <-cq.storage
	}
	return res
}

// DequeueBack iterates over the queue, removing the last element and returning it
func (cq *ChannelQueue[T]) DequeueBack() T {
	var values []T
	storageLen := len(cq.storage)
	for x := 0; x < storageLen; x++ {
		v := <-cq.storage
		values = append(values, v)
	}
	var output T
	for index, v := range values {
		if index == len(values)-1 {
			output = v
		} else {
			cq.storage <- v
		}
	}
	return output
}

// Peek returns (but does not remove) the first element of the queue.
func (cq *ChannelQueue[T]) Peek() T {
	var res T
	if len(cq.storage) == 0 {
		return res
	}
	return cq.Contents()[0]
}

// PeekBack returns (but does not remove) the last element of the queue.
func (cq *ChannelQueue[T]) PeekBack() T {
	var res T
	if len(cq.storage) == 0 {
		return res
	}
	return cq.Contents()[len(cq.storage)-1]
}

// Clear clears the queue.
func (cq *ChannelQueue[T]) Clear() {
	cq.storage = make(chan T, cq.Capacity)
}

// Each pulls every value out of the channel, calls consumer on it, and puts it back.
func (cq *ChannelQueue[T]) Each(consumer func(value T)) {
	if len(cq.storage) == 0 {
		return
	}
	var values []T
	for len(cq.storage) != 0 {
		v := <-cq.storage
		consumer(v)
		values = append(values, v)
	}
	for _, v := range values {
		cq.storage <- v
	}
}

// Consume pulls every value out of the channel, calls consumer on it, effectively clearing the queue.
func (cq *ChannelQueue[T]) Consume(consumer func(value T)) {
	if len(cq.storage) == 0 {
		return
	}
	for len(cq.storage) != 0 {
		v := <-cq.storage
		consumer(v)
	}
}

// EachUntil pulls every value out of the channel, calls consumer on it, and puts it back and can abort mid-process.
func (cq *ChannelQueue[T]) EachUntil(consumer func(value T) bool) {
	contents := cq.Contents()
	for x := 0; x < len(contents); x++ {
		if consumer(contents[x]) {
			return
		}
	}
}

// ReverseEachUntil pulls every value out of the channel, calls consumer on it, and puts it back and can abort mid-process.
func (cq *ChannelQueue[T]) ReverseEachUntil(consumer func(value T) bool) {
	contents := cq.Contents()
	for x := len(contents) - 1; x >= 0; x-- {
		if consumer(contents[x]) {
			return
		}
	}
}

// Contents iterates over the queue and returns a slice of its contents.
func (cq *ChannelQueue[T]) Contents() []T {
	var values []T
	storageLen := len(cq.storage)
	for x := 0; x < storageLen; x++ {
		v := <-cq.storage
		values = append(values, v)
	}
	for _, v := range values {
		cq.storage <- v
	}
	return values
}

// Drain iterates over the queue and returns a slice of its contents, leaving it empty.
func (cq *ChannelQueue[T]) Drain() []T {
	var values []T
	for len(cq.storage) != 0 {
		v := <-cq.storage
		values = append(values, v)
	}
	return values
}
