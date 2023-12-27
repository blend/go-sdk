/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

// NewLinkedList returns a new linked list instance.
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// NewLinkedListFromValues creates a linked list out of a slice.
func NewLinkedListFromValues[T any](values []T) *LinkedList[T] {
	list := new(LinkedList[T])
	for _, v := range values {
		list.Enqueue(v)
	}
	return list
}

// LinkedList is an implementation of a fifo buffer using nodes and pointers.
// Remarks; it is not thread safe. It is constant(ish) time in all ops.
type LinkedList[T any] struct {
	head	*listNode[T]
	tail	*listNode[T]
	length	int
}

// Len returns the length of the queue in constant time.
func (q *LinkedList[T]) Len() int {
	return q.length
}

// Enqueue adds a new value to the queue.
func (q *LinkedList[T]) Enqueue(value T) {
	node := &listNode[T]{Value: value}

	if q.head == nil {	// the queue is empty, that is to say head is nil
		q.head = node
		q.tail = node
	} else {	// the queue is not empty, we have a (valid) tail pointer
		q.tail.Previous = node
		node.Next = q.tail
		q.tail = node
	}

	q.length++
}

// Dequeue removes an item from the front of the queue and returns it.
func (q *LinkedList[T]) Dequeue() T {
	var res T
	if q.head == nil {
		return res
	}

	headValue := q.head.Value

	if q.length == 1 && q.head == q.tail {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.Previous
		if q.head != nil {
			q.head.Next = nil
		}
	}

	q.length--
	return headValue
}

// DequeueBack pops the _last_ element off the linked list.
func (q *LinkedList[T]) DequeueBack() T {
	var res T
	if q.tail == nil {
		return res
	}
	tailValue := q.tail.Value

	if q.length == 1 {
		q.head = nil
		q.tail = nil
	} else {
		q.tail = q.tail.Next
		if q.tail != nil {
			q.tail.Previous = nil
		}
	}

	q.length--
	return tailValue
}

// Peek returns the first element of the queue but does not remove it.
func (q *LinkedList[T]) Peek() T {
	var res T
	if q.head == nil {
		return res
	}
	return q.head.Value
}

// PeekBack returns the last element of the queue.
func (q *LinkedList[T]) PeekBack() T {
	var res T
	if q.tail == nil {
		return res
	}
	return q.tail.Value
}

// Clear clears the linked list.
func (q *LinkedList[T]) Clear() {
	q.tail = nil
	q.head = nil
	q.length = 0
}

// Drain calls the consumer for each element of the linked list.
func (q *LinkedList[T]) Drain() []T {
	if q.head == nil {
		return nil
	}

	contents := make([]T, q.length)
	nodePtr := q.head
	index := 0
	for nodePtr != nil {
		contents[index] = nodePtr.Value
		nodePtr = nodePtr.Previous
		index++
	}
	q.tail = nil
	q.head = nil
	q.length = 0
	return contents
}

// Each calls the consumer for each element of the linked list.
func (q *LinkedList[T]) Each(consumer func(value T)) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		consumer(nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
}

// Consume calls the consumer for each element of the linked list, removing it.
func (q *LinkedList[T]) Consume(consumer func(value T)) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		consumer(nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
	q.tail = nil
	q.head = nil
	q.length = 0
}

// EachUntil calls the consumer for each element of the linked list, but can abort.
func (q *LinkedList[T]) EachUntil(consumer func(value T) bool) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		if !consumer(nodePtr.Value) {
			return
		}
		nodePtr = nodePtr.Previous
	}
}

// ReverseEachUntil calls the consumer for each element of the linked list, but can abort.
func (q *LinkedList[T]) ReverseEachUntil(consumer func(value T) bool) {
	if q.head == nil {
		return
	}

	nodePtr := q.tail
	for nodePtr != nil {
		if !consumer(nodePtr.Value) {
			return
		}
		nodePtr = nodePtr.Next
	}
}

// Contents returns the full contents of the queue as a slice.
func (q *LinkedList[T]) Contents() []T {
	if q.head == nil {
		return nil
	}

	var values []T
	nodePtr := q.head
	for nodePtr != nil {
		values = append(values, nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
	return values
}

type listNode[T any] struct {
	Next		*listNode[T]
	Previous	*listNode[T]
	Value		T
}
