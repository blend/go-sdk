/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

import (
	"testing"
	"time"
)

const (
	// DefaultSampleSize is the default number of steps to run per test.
	DefaultSampleSize	= 10000

	// DefaultStasisSize is the stasis size for the fixed length test.
	DefaultStasisSize	= 512
)

// QueueFactory is a function that emits a Queue[time.Time]
type QueueFactory func(capacity int) Queue[time.Time]

func doQueueBenchmark(queueFactory QueueFactory, sampleSize int, b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		q := queueFactory(sampleSize)
		for x := 0; x < sampleSize; x++ {
			q.Enqueue(time.Now().UTC())
		}
		for x := 0; x < sampleSize; x++ {
			q.Dequeue()
		}
	}
}

func doFixedQueueBenchmark(queueFactory QueueFactory, sampleSize, stasisSize int, b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		q := queueFactory(stasisSize)
		for x := 0; x < sampleSize; x++ {
			q.Enqueue(time.Now().UTC())
			if q.Len() < stasisSize {
				continue
			}
			q.Dequeue()
		}
	}
}

func makeLinkedList[T any](capacity int) Queue[T] {
	return NewLinkedList[T]()
}

func makeChannelQueue[T any](capacity int) Queue[T] {
	return NewChannelQueueWithCapacity[T](capacity)
}

func makeRingBuffer[T any](capacity int) Queue[T] {
	return NewRingBufferWithCapacity[T](capacity)
}

func makeSyncedRingBuffer[T any](capacity int) Queue[T] {
	rb := NewSyncRingBufferWithCapacity[T](capacity)
	return rb
}

func BenchmarkLinkedList(b *testing.B) {
	doQueueBenchmark(makeLinkedList[time.Time], DefaultSampleSize, b)
}

func BenchmarkChannelQueue(b *testing.B) {
	doQueueBenchmark(makeChannelQueue[time.Time], DefaultSampleSize, b)
}

func BenchmarkRingBuffer(b *testing.B) {
	doQueueBenchmark(makeRingBuffer[time.Time], DefaultSampleSize, b)
}

func BenchmarkRingBufferSynced(b *testing.B) {
	doQueueBenchmark(makeSyncedRingBuffer[time.Time], DefaultSampleSize, b)
}

func BenchmarkFixedLinkedList(b *testing.B) {
	doFixedQueueBenchmark(makeLinkedList[time.Time], DefaultSampleSize, DefaultStasisSize, b)
}

func BenchmarkFixedChannelQueue(b *testing.B) {
	doFixedQueueBenchmark(makeChannelQueue[time.Time], DefaultSampleSize, DefaultStasisSize, b)
}

func BenchmarkFixedRingBuffer(b *testing.B) {
	doFixedQueueBenchmark(makeRingBuffer[time.Time], DefaultSampleSize, DefaultStasisSize, b)
}
