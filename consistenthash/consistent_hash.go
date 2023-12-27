/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package consistenthash

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

const (
	// DefaultReplicas is the default number of bucket virtual replicas.
	DefaultReplicas = 16
)

var (
	_ json.Marshaler = (*ConsistentHash)(nil)
	_ fmt.Stringer   = (*ConsistentHash)(nil)
)

// New returns a new consistent hash.
func New(opts ...Option) *ConsistentHash {
	var ch ConsistentHash
	for _, opt := range opts {
		opt(&ch)
	}
	return &ch
}

// Option mutates a consistent hash.
type Option func(*ConsistentHash)

// OptBuckets adds buckets to the consistent hash.
//
// It is functionally equiavalent to looping over the buckets
// and calling `AddBuckets(bucketsj...)` for it.
func OptBuckets(buckets ...string) Option {
	return func(ch *ConsistentHash) {
		ch.AddBuckets(buckets...)
	}
}

// OptReplicas sets the bucket virtual replica count.
//
// More virtual replicas can help with making item assignments
// more uniform, but the tradeoff is every operation takes a little
// longer as log2 of the number of buckets times the number of virtual replicas.
//
// If not provided, the default (16) is used.
func OptReplicas(replicas int) Option {
	return func(ch *ConsistentHash) { ch.replicas = replicas }
}

// OptHashFunction sets the hash function.
//
// The default hash function is `consistenthash.StableHash` which uses
// a stable crc64 hash function to preserve ordering between process restarts.
func OptHashFunction(hashFunction HashFunction) Option {
	return func(ch *ConsistentHash) { ch.hashFunction = hashFunction }
}

// ConsistentHash creates hashed assignments for each bucket.
type ConsistentHash struct {
	mu sync.RWMutex

	replicas     int
	buckets      map[string]struct{}
	hashFunction HashFunction
	hashring     []HashedBucket
}

//
// properties with defaults
//

// ReplicasOrDefault is the default number of bucket virtual replicas.
func (ch *ConsistentHash) ReplicasOrDefault() int {
	if ch.replicas > 0 {
		return ch.replicas
	}
	return DefaultReplicas
}

// HashFunctionOrDefault returns the provided hash function or a default.
func (ch *ConsistentHash) HashFunctionOrDefault() HashFunction {
	if ch.hashFunction != nil {
		return ch.hashFunction
	}
	return StableHash
}

//
// Write methods
//

// AddBuckets adds a list of buckets to the consistent hash, and returns
// a boolean indiciating if _any_ buckets were added.
//
// If any of the new buckets do not exist on the hash ring the
// new bucket will be inserted `ReplicasOrDefault` number
// of times into the internal hashring.
//
// If any of the new buckets already exist on the hash ring
// no action is taken for that bucket.
//
// Calling `AddBuckets` is safe to do concurrently
// and acquires a write lock on the consistent hash reference.
func (ch *ConsistentHash) AddBuckets(newBuckets ...string) (ok bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.buckets == nil {
		ch.buckets = make(map[string]struct{})
	}
	for _, newBucket := range newBuckets {
		if _, ok := ch.buckets[newBucket]; ok {
			continue
		}
		ok = true
		ch.buckets[newBucket] = struct{}{}
		ch.insertUnsafe(newBucket)
	}
	return
}

// RemoveBucket removes a bucket from the consistent hash, and returns
// a boolean indicating if the provided bucket was found.
//
// If the bucket exists on the hash ring, the bucket and its replicas are removed.
//
// If the bucket does not exist on the ring, no action is taken.
//
// Calling `RemoveBucket` is safe to do concurrently
// and acquires a write lock on the consistent hash reference.
func (ch *ConsistentHash) RemoveBucket(toRemove string) (ok bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.buckets == nil {
		return
	}
	if _, ok = ch.buckets[toRemove]; !ok {
		return
	}
	delete(ch.buckets, toRemove)
	for x := 0; x < ch.ReplicasOrDefault(); x++ {
		index := ch.search(ch.bucketHashKey(toRemove, x))
		ch.hashring = append(ch.hashring[:index], ch.hashring[index+1:]...)
	}
	return
}

//
// Read methods
//

// Buckets returns the buckets.
//
// Calling `Buckets` is safe to do concurrently and acquires
// a read lock on the consistent hash reference.
func (ch *ConsistentHash) Buckets() (buckets []string) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	for bucket := range ch.buckets {
		buckets = append(buckets, bucket)
	}
	sort.Strings(buckets)
	return
}

// Assignment returns the bucket assignment for a given item.
//
// Calling `Assignment` is safe to do concurrently and acquires
// a read lock on the consistent hash reference.
func (ch *ConsistentHash) Assignment(item string) (bucket string) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	bucket = ch.assignmentUnsafe(item)
	return
}

// IsAssigned returns if a given bucket is assigned a given item.
//
// Calling `IsAssigned` is safe to do concurrently and acquires
// a read lock on the consistent hash reference.
func (ch *ConsistentHash) IsAssigned(bucket, item string) (ok bool) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	ok = bucket == ch.assignmentUnsafe(item)
	return
}

// Assignments returns the assignments for a given list of items organized
// by the name of the bucket, and an array of the assigned items.
//
// Calling `Assignments` is safe to do concurrently and acquires
// a read lock on the consistent hash reference.
func (ch *ConsistentHash) Assignments(items ...string) map[string][]string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	output := make(map[string][]string)
	for _, item := range items {
		bucket := ch.assignmentUnsafe(item)
		output[bucket] = append(output[bucket], item)
	}
	return output
}

// String returns a string form of the hash for debugging purposes.
//
// Calling `String` is safe to do concurrently and acquires
// a read lock on the consistent hash reference.
func (ch *ConsistentHash) String() string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	var output []string
	for _, bucket := range ch.hashring {
		output = append(output, fmt.Sprintf("%d:%s-%02d", bucket.Hashcode, bucket.Bucket, bucket.Replica))
	}
	return strings.Join(output, ", ")
}

// MarshalJSON marshals the consistent hash as json.
//
// The form of the returned json is the underlying []HashedBucket
// and there is no corresponding `UnmarshalJSON` because
// it is uncertain on the other end what the hashfunction is
// because functions can't be json serialized.
//
// You should use MarshalJSON for communicating information
// for debugging purposes only.
//
// Calling `MarshalJSON` is safe to do concurrently and acquires
// a read lock on the consistent hash reference.
func (ch *ConsistentHash) MarshalJSON() ([]byte, error) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	return json.Marshal(ch.hashring)
}

//
// internal / unexported helpers
//

// assignmentUnsafe searches for the item's matching bucket based
// on a binary search, and if the index returned is outside the
// ring length, the first index (0) is returned to simulate wrapping around.
func (ch *ConsistentHash) assignmentUnsafe(item string) (bucket string) {
	index := ch.search(item)
	if index >= len(ch.hashring) {
		index = 0
	}
	bucket = ch.hashring[index].Bucket
	return
}

// insert inserts a hashring bucket.
//
// insert uses an insertion sort such that the
// resulting ring will remain sorted after insert.
//
// it will also insert `ReplicasOrDefault` copies of the bucket
// to help distribute items across buckets more evenly.
func (ch *ConsistentHash) insertUnsafe(bucket string) {
	for x := 0; x < ch.ReplicasOrDefault(); x++ {
		ch.hashring = InsertionSort(ch.hashring, HashedBucket{
			Bucket:   bucket,
			Replica:  x,
			Hashcode: ch.hashcode(ch.bucketHashKey(bucket, x)),
		})
	}
}

// search does a binary search for the first hashring index whose
// node hashcode is >= the hashcode of a given item.
func (ch *ConsistentHash) search(item string) (index int) {
	return sort.Search(len(ch.hashring), ch.searchFn(ch.hashcode(item)))
}

// searchFn returns a closure searching for a given hashcode.
func (ch *ConsistentHash) searchFn(hashcode uint64) func(index int) bool {
	return func(index int) bool {
		return ch.hashring[index].Hashcode >= hashcode
	}
}

// bucketHashKey formats a hash key for a given bucket virtual replica.
func (ch *ConsistentHash) bucketHashKey(bucket string, index int) string {
	return bucket + "|" + fmt.Sprintf("%02d", index)
}

// hashcode creates a hashcode for a given string
func (ch *ConsistentHash) hashcode(item string) uint64 {
	return ch.HashFunctionOrDefault()([]byte(item))
}

// HashedBucket is a bucket in the hashring
// that holds the hashcode, the bucket name (as Bucket)
// and the virtual replica index.
type HashedBucket struct {
	Hashcode uint64 `json:"hashcode"`
	Bucket   string `json:"bucket"`
	Replica  int    `json:"replica"`
}
