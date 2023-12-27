/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

// Queue is an interface for implementations of a FIFO buffer.
type Queue[T any] interface {
	Len() int
	Enqueue(T)
	Dequeue() T
	DequeueBack() T
	Peek() T
	PeekBack() T
	Drain() []T
	Contents() []T
	Clear()

	Consume(consumer func(value T))
	Each(consumer func(value T))
	EachUntil(consumer func(value T) bool)
	ReverseEachUntil(consumer func(value T) bool)
}
