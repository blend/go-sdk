/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

import (
	"fmt"
	"strings"
)

const (
	ringBufferMinimumGrow     = 4
	ringBufferGrowFactor      = 200
	ringBufferDefaultCapacity = 4
)

// NewRingBuffer creates a new, empty, RingBuffer.
func NewRingBuffer[T any]() *RingBuffer[T] {
	return &RingBuffer[T]{
		slice: make([]T, ringBufferDefaultCapacity),
		head:  0,
		tail:  0,
		size:  0,
	}
}

// NewRingBufferWithCapacity creates a new ring buffer with a given capacity.
func NewRingBufferWithCapacity[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		slice: make([]T, capacity),
		head:  0,
		tail:  0,
		size:  0,
	}
}

// NewRingBufferFromValues creates a new ring buffer out of a slice.
func NewRingBufferFromValues[T any](values []T) *RingBuffer[T] {
	return &RingBuffer[T]{
		slice: values,
		head:  0,
		tail:  len(values) - 1,
		size:  len(values),
	}
}

// RingBuffer is a fifo buffer that is backed by a pre-allocated slice, instead of allocating
// a whole new node object for each element (which saves GC churn).
// Enqueue can be O(n), Dequeue can be O(1).
type RingBuffer[T any] struct {
	slice []T
	head  int
	tail  int
	size  int
}

// Len returns the length of the ring buffer (as it is currently populated).
// Actual memory footprint may be different.
func (rb *RingBuffer[T]) Len() (len int) {
	return rb.size
}

// Capacity returns the total size of the ring buffer, including empty elements.
func (rb *RingBuffer[T]) Capacity() int {
	return len(rb.slice)
}

// Clear removes all objects from the RingBuffer.
func (rb *RingBuffer[T]) Clear() {
	if rb.head < rb.tail {
		sliceClear(rb.slice, rb.head, rb.size)
	} else {
		sliceClear(rb.slice, rb.head, len(rb.slice)-rb.head)
		sliceClear(rb.slice, 0, rb.tail)
	}

	rb.head = 0
	rb.tail = 0
	rb.size = 0
}

// Enqueue adds an element to the "back" of the RingBuffer.
func (rb *RingBuffer[T]) Enqueue(value T) {
	if rb.size == len(rb.slice) {
		newCapacity := int(len(rb.slice) * int(ringBufferGrowFactor/100))
		if newCapacity < (len(rb.slice) + ringBufferMinimumGrow) {
			newCapacity = len(rb.slice) + ringBufferMinimumGrow
		}
		rb.setCapacity(newCapacity)
	}

	rb.slice[rb.tail] = value
	rb.tail = (rb.tail + 1) % len(rb.slice)
	rb.size++
}

// Dequeue removes the first (oldest) element from the RingBuffer.
func (rb *RingBuffer[T]) Dequeue() T {
	var res T
	if rb.size == 0 {
		return res
	}

	removed := rb.slice[rb.head]
	rb.head = (rb.head + 1) % len(rb.slice)
	rb.size--

	return removed
}

// DequeueBack removes the last (newest) element from the RingBuffer.
func (rb *RingBuffer[T]) DequeueBack() T {
	var res T
	if rb.size == 0 {
		return res
	}

	// tail is the
	var removed T
	if rb.tail == 0 {
		removed = rb.slice[len(rb.slice)-1]
		rb.tail = len(rb.slice) - 1
	} else {
		removed = rb.slice[rb.tail-1]
		rb.tail = rb.tail - 1
	}
	rb.size--
	return removed
}

// Peek returns but does not remove the first element.
func (rb *RingBuffer[T]) Peek() T {
	var res T
	if rb.size == 0 {
		return res
	}
	return rb.slice[rb.head]
}

// PeekBack returns but does not remove the last element.
func (rb *RingBuffer[T]) PeekBack() T {
	var res T
	if rb.size == 0 {
		return res
	}
	if rb.tail == 0 {
		return rb.slice[len(rb.slice)-1]
	}
	return rb.slice[rb.tail-1]
}

func (rb *RingBuffer[T]) setCapacity(capacity int) {
	newSlice := make([]T, capacity)
	if rb.size > 0 {
		if rb.head < rb.tail {
			sliceCopy(rb.slice, rb.head, newSlice, 0, rb.size)
		} else {
			sliceCopy(rb.slice, rb.head, newSlice, 0, len(rb.slice)-rb.head)
			sliceCopy(rb.slice, 0, newSlice, len(rb.slice)-rb.head, rb.tail)
		}
	}
	rb.slice = newSlice
	rb.head = 0
	rb.tail = 0
	if rb.size != capacity {
		rb.tail = rb.size
	}
}

// trimExcess resizes the buffer to better fit the contents.
func (rb *RingBuffer[T]) trimExcess() {
	threshold := float64(len(rb.slice)) * 0.9
	if rb.size < int(threshold) {
		rb.setCapacity(rb.size)
	}
}

// Contents returns the ring buffer, in order, as a slice.
func (rb *RingBuffer[T]) Contents() []T {
	newSlice := make([]T, rb.size)

	if rb.size == 0 {
		return newSlice
	}

	if rb.head < rb.tail {
		sliceCopy(rb.slice, rb.head, newSlice, 0, rb.size)
		sliceClear(rb.slice, rb.head, rb.size)
	} else {
		sliceCopy(rb.slice, rb.head, newSlice, 0, len(rb.slice)-rb.head)
		sliceClear(rb.slice, rb.head, len(rb.slice)-rb.head)
		sliceCopy(rb.slice, 0, newSlice, len(rb.slice)-rb.head, rb.tail)
		sliceClear(rb.slice, 0, rb.tail)
	}

	return newSlice
}

// Drain clears the buffer and removes the contents.
func (rb *RingBuffer[T]) Drain() []T {
	newSlice := make([]T, rb.size)

	if rb.size == 0 {
		return newSlice
	}

	if rb.head < rb.tail {
		sliceCopy(rb.slice, rb.head, newSlice, 0, rb.size)
	} else {
		sliceCopy(rb.slice, rb.head, newSlice, 0, len(rb.slice)-rb.head)
		sliceCopy(rb.slice, 0, newSlice, len(rb.slice)-rb.head, rb.tail)
	}

	rb.head = 0
	rb.tail = 0
	rb.size = 0

	return newSlice
}

// Each calls the consumer for each element in the buffer.
func (rb *RingBuffer[T]) Each(consumer func(value T)) {
	if rb.size == 0 {
		return
	}

	if rb.head < rb.tail {
		for cursor := rb.head; cursor < rb.tail; cursor++ {
			consumer(rb.slice[cursor])
		}
	} else {
		for cursor := rb.head; cursor < len(rb.slice); cursor++ {
			consumer(rb.slice[cursor])
		}
		for cursor := 0; cursor < rb.tail; cursor++ {
			consumer(rb.slice[cursor])
		}
	}
}

// Consume calls the consumer for each element in the buffer, while also dequeueing that entry.
func (rb *RingBuffer[T]) Consume(consumer func(value T)) {
	if rb.size == 0 {
		return
	}

	length := rb.Len()
	for i := 0; i < length; i++ {
		consumer(rb.Dequeue())
	}
}

// EachUntil calls the consumer for each element in the buffer with a stopping condition in head=>tail order.
func (rb *RingBuffer[T]) EachUntil(consumer func(value T) bool) {
	if rb.size == 0 {
		return
	}

	if rb.head < rb.tail {
		for cursor := rb.head; cursor < rb.tail; cursor++ {
			if !consumer(rb.slice[cursor]) {
				return
			}
		}
	} else {
		for cursor := rb.head; cursor < len(rb.slice); cursor++ {
			if !consumer(rb.slice[cursor]) {
				return
			}
		}
		for cursor := 0; cursor < rb.tail; cursor++ {
			if !consumer(rb.slice[cursor]) {
				return
			}
		}
	}
}

// ReverseEachUntil calls the consumer for each element in the buffer with a stopping condition in tail=>head order.
func (rb *RingBuffer[T]) ReverseEachUntil(consumer func(value T) bool) {
	if rb.size == 0 {
		return
	}

	if rb.head < rb.tail {
		for cursor := rb.tail - 1; cursor >= rb.head; cursor-- {
			if !consumer(rb.slice[cursor]) {
				return
			}
		}
	} else {
		for cursor := rb.tail; cursor > 0; cursor-- {
			if !consumer(rb.slice[cursor]) {
				return
			}
		}
		for cursor := len(rb.slice) - 1; cursor >= rb.head; cursor-- {
			if !consumer(rb.slice[cursor]) {
				return
			}
		}
	}
}

func (rb *RingBuffer[T]) String() string {
	var values []string
	for _, elem := range rb.Contents() {
		values = append(values, fmt.Sprintf("%v", elem))
	}
	return strings.Join(values, " <= ")
}

func sliceClear[T any](source []T, index, length int) {
	var val T
	for x := 0; x < length; x++ {
		absoluteIndex := x + index
		source[absoluteIndex] = val
	}
}

func sliceCopy[T any](source []T, sourceIndex int, destination []T, destinationIndex, length int) {
	for x := 0; x < length; x++ {
		from := sourceIndex + x
		to := destinationIndex + x

		destination[to] = source[from]
	}
}
