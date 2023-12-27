/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_LinkedList(t *testing.T) {
	a := assert.New(t)

	q := NewLinkedList[string]()
	a.Empty(q.head)
	a.Empty(q.tail)
	a.Empty(q.Contents())
	a.Empty(q.Dequeue())
	a.Equal(q.head, q.tail)
	a.Empty(q.Peek())
	a.Empty(q.PeekBack())
	a.Equal(0, q.Len())

	q.Enqueue("foo")
	a.NotNil(q.head)
	a.Empty(q.head.Previous)
	a.NotNil(q.tail)
	a.Empty(q.tail.Previous)
	a.Equal(q.head, q.tail)
	a.Equal(1, q.Len())
	a.Equal("foo", q.Peek())
	a.Equal("foo", q.PeekBack())

	q.Enqueue("bar")
	a.NotNil(q.head)
	a.NotNil(q.head.Previous)
	a.Empty(q.head.Previous.Previous)
	a.Equal(q.head.Previous, q.tail)
	a.NotNil(q.tail)
	a.NotNil(q.tail.Next)
	a.NotEqual(q.head, q.tail)
	a.Equal(2, q.Len())
	a.Equal("foo", q.Peek())
	a.Equal("bar", q.PeekBack())

	q.Enqueue("baz")
	a.NotNil(q.head)
	a.NotNil(q.head.Previous)
	a.NotNil(q.head.Previous.Previous)
	a.Empty(q.head.Previous.Previous.Previous)
	a.Equal(q.head.Previous.Previous, q.tail)
	a.NotNil(q.tail)
	a.NotEqual(q.head, q.tail)
	a.Equal(3, q.Len())
	a.Equal("foo", q.Peek())
	a.Equal("baz", q.PeekBack())

	q.Enqueue("fizz")
	a.NotNil(q.head)
	a.NotNil(q.head.Previous)
	a.NotNil(q.head.Previous.Previous)
	a.NotNil(q.head.Previous.Previous.Previous)
	a.Empty(q.head.Previous.Previous.Previous.Previous)
	a.Equal(q.head.Previous.Previous.Previous, q.tail)
	a.NotNil(q.tail)
	a.NotEqual(q.head, q.tail)
	a.Equal(4, q.Len())
	a.Equal("foo", q.Peek())
	a.Equal("fizz", q.PeekBack())

	values := q.Contents()
	a.Len(values, 4)
	a.Equal("foo", values[0])
	a.Equal("bar", values[1])
	a.Equal("baz", values[2])
	a.Equal("fizz", values[3])

	shouldBeFoo := q.Dequeue()
	a.Equal("foo", shouldBeFoo)
	a.NotNil(q.head)
	a.NotNil(q.head.Previous)
	a.NotNil(q.head.Previous.Previous)
	a.Empty(q.head.Previous.Previous.Previous)
	a.Equal(q.head.Previous.Previous, q.tail)
	a.NotNil(q.tail)
	a.NotEqual(q.head, q.tail)
	a.Equal(q.Len(), 3)
	a.Equal("bar", q.Peek())
	a.Equal("fizz", q.PeekBack())

	shouldBeBar := q.Dequeue()
	a.Equal("bar", shouldBeBar)
	a.NotNil(q.head)
	a.NotNil(q.head.Previous)
	a.Empty(q.head.Previous.Previous)
	a.Equal(q.head.Previous, q.tail)
	a.NotNil(q.tail)
	a.NotEqual(q.head, q.tail)
	a.Equal(2, q.Len())
	a.Equal("baz", q.Peek())
	a.Equal("fizz", q.PeekBack())

	shouldBeBaz := q.Dequeue()
	a.Equal("baz", shouldBeBaz)
	a.NotNil(q.head)
	a.Empty(q.head.Previous)
	a.NotNil(q.tail)
	a.Equal(q.head, q.tail)
	a.Equal(1, q.Len())
	a.Equal("fizz", q.Peek())
	a.Equal("fizz", q.PeekBack())

	shouldBeFizz := q.Dequeue()
	a.Equal("fizz", shouldBeFizz)
	a.Empty(q.head)
	a.Empty(q.tail)
	a.Empty(q.Peek())
	a.Empty(q.PeekBack())
	a.Equal(0, q.Len())

	q.Enqueue("foo")
	q.Enqueue("bar")
	q.Enqueue("baz")
	a.Equal(3, q.Len())
	q.Clear()
	a.Equal(0, q.Len())

	q.Enqueue("foo")
	q.Enqueue("bar")
	q.Enqueue("baz")

	var items []string
	q.Each(func(v string) {
		items = append(items, v)
	})
	a.Len(items, 3)
	a.Equal("foo", items[0])
	a.Equal("bar", items[1])
	a.Equal("baz", items[2])
	a.Equal(3, q.Len())

	items = []string{}
	q.Consume(func(v string) {
		items = append(items, v)
	})
	a.Equal(0, q.Len())
	a.Len(items, 3)
	a.Equal("foo", items[0])
	a.Equal("bar", items[1])
	a.Equal("baz", items[2])

	q.Enqueue("foo")
	q.Enqueue("bar")
	q.Enqueue("baz")
	a.Equal(3, q.Len())

	contents := q.Drain()
	a.Len(contents, 3)
	a.Equal(0, q.Len())
}

func Test_LinkedList_DequeueBack(t *testing.T) {
	a := assert.New(t)

	q := NewLinkedList[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)

	a.Equal(4, q.DequeueBack())
	a.Equal(3, q.DequeueBack())
	a.Equal(2, q.DequeueBack())
	a.Equal(1, q.DequeueBack())
	a.Empty(q.DequeueBack())
	a.Empty(q.DequeueBack())

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)

	a.Equal(4, q.DequeueBack())
	a.Equal(3, q.DequeueBack())
	a.Equal(2, q.DequeueBack())
	a.Equal(1, q.DequeueBack())
	a.Empty(q.DequeueBack())
}
