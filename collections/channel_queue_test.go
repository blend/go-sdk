package collections

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestConcurrentQueue(t *testing.T) {
	a := assert.New(t)

	q := NewChannelQueueWithCapacity(4)
	a.Empty(q.Contents())
	a.Nil(q.Dequeue())
	a.Equal(0, q.Len())

	q.Enqueue("foo")
	a.Equal(1, q.Len())
	a.Equal("foo", q.Peek())
	a.Equal("foo", q.PeekBack())

	q.Enqueue("bar")
	a.Equal(2, q.Len())

	q.Enqueue("baz")
	a.Equal(3, q.Len())

	q.Enqueue("fizz")
	a.Equal(q.Len(), 4)

	values := q.Contents()
	a.Len(values, 4)
	a.Equal("foo", values[0])
	a.Equal("bar", values[1])
	a.Equal("baz", values[2])
	a.Equal("fizz", values[3])

	shouldBeFoo := q.Dequeue()
	a.Equal("foo", shouldBeFoo)
	a.Equal(q.Len(), 3)

	shouldBeBar := q.Dequeue()
	a.Equal("bar", shouldBeBar)
	a.Equal(2, q.Len())

	shouldBeBaz := q.Dequeue()
	a.Equal("baz", shouldBeBaz)
	a.Equal(1, q.Len())

	shouldBeFizz := q.Dequeue()
	a.Equal("fizz", shouldBeFizz)
	a.Equal(0, q.Len())

	q.Enqueue("foo")
	a.Equal(1, q.Len())
	q.Clear()
	a.Equal(0, q.Len())

	q.Enqueue("foo")
	q.Enqueue("bar")
	q.Enqueue("baz")

	var items []string
	q.Consume(func(v Any) {
		items = append(items, v.(string))
	})
	a.Equal(0, q.Len())
	a.Len(items, 3)
	a.Equal("foo", items[0])
	a.Equal("bar", items[1])
	a.Equal("baz", items[2])

	q.Enqueue("foo")
	q.Enqueue("bar")
	q.Enqueue("baz")

	items = []string{}
	q.Each(func(v Any) {
		items = append(items, v.(string))
	})
	a.Equal(3, q.Len())
	a.Len(items, 3)
	a.Equal("foo", items[0])
	a.Equal("bar", items[1])
	a.Equal("baz", items[2])

	contents := q.Drain()
	a.Len(contents, 3)
}
