/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package consistenthash

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_ConsistentHash_typical(t *testing.T) {
	its := assert.New(t)

	const bucketCount = 5
	const itemCount = 150

	var buckets []string
	for x := 0; x < bucketCount; x++ {
		buckets = append(buckets, fmt.Sprintf("worker-%d", x))
	}
	var items []string
	for x := 0; x < itemCount; x++ {
		items = append(items, fmt.Sprintf("google-%d", x))
	}

	ch := New(
		OptBuckets(
			buckets...,
		),
	)
	its.Len(ch.hashring, bucketCount*ch.ReplicasOrDefault(), "the internal hashring should mirror the count of the buckets")

	returnedBuckets := ch.Buckets()
	its.Len(returnedBuckets, bucketCount)
	setIsSorted(its, returnedBuckets)

	assignments := ch.Assignments(items...)

	worker0items := assignments["worker-0"]
	its.NotEmpty(worker0items)

	worker1items := assignments["worker-1"]
	its.NotEmpty(worker1items)

	worker2items := assignments["worker-2"]
	its.NotEmpty(worker2items)

	worker3items := assignments["worker-3"]
	its.NotEmpty(worker3items)

	worker4items := assignments["worker-4"]
	its.NotEmpty(worker4items)

	// verify that all the bucket assignments are disjoint, that is
	// none of the items in an assignment exist in another assignment
	setsAreDisjoint(its, worker0items, worker1items, worker2items, worker3items, worker4items)

	// verify this is also the case through the `IsAssigned` method
	for _, item := range worker0items {
		its.True(ch.IsAssigned("worker-0", item))
	}
	// verify the worker-0 items are not assigned to
	// any other nodes
	for _, item := range worker0items {
		its.False(ch.IsAssigned("worker-1", item))
		its.False(ch.IsAssigned("worker-2", item))
		its.False(ch.IsAssigned("worker-3", item))
		its.False(ch.IsAssigned("worker-4", item))
	}

	t.Log(spew(assignmentCounts(assignments)))
	// verify the sets are relatively evenly sized
	// this is generally the most likely test to fail, because
	// the hash bucket assignments can be lumpy.
	setsAreEvenlySized(its, itemCount, worker0items, worker1items, worker2items, worker3items, worker4items)

	// verify that all the consistent hash items exist in exactly one bucket assignment
	setExistsInOtherSets(its, items, worker0items, worker1items, worker2items, worker3items, worker4items)
}

func Test_ConsistentHash_AddBuckets(t *testing.T) {
	its := assert.New(t)

	ch := New()

	res := ch.AddBuckets("worker-0", "worker-1", "worker-2")
	its.True(res)
	res = ch.AddBuckets("worker-0", "worker-1", "worker-2")
	its.False(res, "we should return false if _no_ new buckets were added")
	res = ch.AddBuckets("worker-0", "worker-1", "worker-2", "worker-3")
	its.True(res, "we should return true if _any_ new buckets were added")
	buckets := ch.Buckets()
	its.Len(buckets, 4)
	its.Equal([]string{"worker-0", "worker-1", "worker-2", "worker-3"}, buckets)
}

func Test_ConsistentHash_RemoveBucket(t *testing.T) {
	its := assert.New(t)

	ch := New()

	res := ch.AddBuckets("worker-0", "worker-1", "worker-2")
	its.True(res)
	res = ch.RemoveBucket("worker-3")
	its.False(res, "we should return false if bucket not found")
	res = ch.RemoveBucket("worker-2")
	its.True(res, "we should return true if bucket found")
	buckets := ch.Buckets()
	its.Len(buckets, 2)
	its.Equal([]string{"worker-0", "worker-1"}, buckets)
}

func Test_ConsistentHash_redistribute_addBuckets(t *testing.T) {
	its := assert.New(t)

	const bucketCount = 5
	const itemCount = 100
	const itemsPerBucket = itemCount / bucketCount
	const maxBucketDelta = (itemsPerBucket / bucketCount) + 4

	var buckets []string
	for x := 0; x < bucketCount; x++ {
		buckets = append(buckets, fmt.Sprintf("worker-%d", x))
	}
	var items []string
	for x := 0; x < itemCount; x++ {
		items = append(items, fmt.Sprintf("google-%d", x))
	}

	ch := New(
		OptBuckets(
			buckets...,
		),
	)

	assignments := ch.Assignments(items...)

	oldWorker0items := assignments["worker-0"]
	its.NotEmpty(oldWorker0items)

	oldWorker1items := assignments["worker-1"]
	its.NotEmpty(oldWorker1items)

	oldWorker2items := assignments["worker-2"]
	its.NotEmpty(oldWorker2items)

	oldWorker3items := assignments["worker-3"]
	its.NotEmpty(oldWorker3items)

	oldWorker4items := assignments["worker-4"]
	its.NotEmpty(oldWorker4items)

	// simulate adding a bucket
	ch.AddBuckets("worker-5")
	its.Len(ch.buckets, bucketCount+1)
	its.Len(ch.hashring, (bucketCount+1)*ch.ReplicasOrDefault())

	newAssignments := ch.Assignments(items...)
	its.Len(newAssignments, bucketCount+1, "assignments length should mirror buckets")

	worker0items := newAssignments["worker-0"]
	its.NotEmpty(worker0items)

	worker1items := newAssignments["worker-1"]
	its.NotEmpty(worker1items)

	worker2items := newAssignments["worker-2"]
	its.NotEmpty(worker2items)

	worker3items := newAssignments["worker-3"]
	its.NotEmpty(worker3items)

	worker4items := newAssignments["worker-4"]
	its.NotEmpty(worker4items)

	worker5items := newAssignments["worker-5"]
	its.NotEmpty(worker5items)

	// verify that all the bucket assignments are disjoint, that is
	// none of the items in an assignment exist in another assignment
	setsAreDisjoint(its, worker0items, worker1items, worker2items, worker3items, worker4items, worker5items)

	// verify that all the consistent hash items exist in a bucket assignment
	setExistsInOtherSets(its, items, worker0items, worker1items, worker2items, worker3items, worker4items, worker5items)

	// verify that we're moving items around consistently

	t.Log(spew(assignmentCounts(newAssignments)))
	itsConsistent(its, oldWorker0items, worker0items, maxBucketDelta)
	itsConsistent(its, oldWorker1items, worker1items, maxBucketDelta)
	itsConsistent(its, oldWorker2items, worker2items, maxBucketDelta)
	itsConsistent(its, oldWorker3items, worker3items, maxBucketDelta)
	itsConsistent(its, oldWorker4items, worker4items, maxBucketDelta)
}

func Test_ConsistentHash_redistribute_removeBucket(t *testing.T) {
	its := assert.New(t)

	const bucketCount = 5
	const itemCount = 100
	const itemsPerBucket = itemCount / bucketCount
	const maxBucketDelta = (itemsPerBucket / bucketCount) + 4

	var buckets []string
	for x := 0; x < bucketCount; x++ {
		buckets = append(buckets, fmt.Sprintf("worker-%d", x))
	}
	var items []string
	for x := 0; x < itemCount; x++ {
		items = append(items, fmt.Sprintf("google-%d", x))
	}

	ch := New(
		OptBuckets(
			buckets...,
		),
	)
	assignments := ch.Assignments(items...)

	oldWorker0items := assignments["worker-0"]
	its.NotEmpty(oldWorker0items)

	oldWorker1items := assignments["worker-1"]
	its.NotEmpty(oldWorker1items)

	oldWorker2items := assignments["worker-2"]
	its.NotEmpty(oldWorker2items)

	oldWorker3items := assignments["worker-3"]
	its.NotEmpty(oldWorker3items)

	oldWorker4items := assignments["worker-4"]
	its.NotEmpty(oldWorker4items)

	// the maximum number of items to move around is the
	// single bucket we're removing (plus a fudge factor)

	// simulate dropping a bucket (or node)
	its.True(ch.RemoveBucket("worker-2"))
	its.Len(ch.buckets, bucketCount-1)
	its.Len(ch.hashring, (bucketCount-1)*ch.ReplicasOrDefault())

	_, ok := ch.buckets["worker-2"]
	its.False(ok)
	for _, bucket := range ch.hashring {
		its.NotEqual("worker-2", bucket.Bucket)
	}

	assignments = ch.Assignments(items...)
	its.Len(assignments, len(ch.buckets), "assignments length should mirror buckets")

	worker0items := assignments["worker-0"]
	its.NotEmpty(oldWorker0items)

	worker1items := assignments["worker-1"]
	its.NotEmpty(oldWorker1items)

	worker2items := assignments["worker-2"]
	its.NotEmpty(oldWorker2items)

	worker3items := assignments["worker-3"]
	its.NotEmpty(oldWorker3items)

	worker4items := assignments["worker-4"]
	its.NotEmpty(oldWorker4items)

	// verify that all the bucket assignments are disjoint, that is
	// none of the items in an assignment exist in another assignment
	setsAreDisjoint(its, worker0items, worker1items, worker2items, worker3items, worker4items)

	// verify that all the consistent hash items exist in a bucket assignment
	setExistsInOtherSets(its, items, worker0items, worker1items, worker3items, worker4items)

	// verify that we're moving items around consistently
	itsConsistent(its, oldWorker0items, worker0items, maxBucketDelta)
	itsConsistent(its, oldWorker1items, worker1items, maxBucketDelta)
	itsConsistent(its, oldWorker3items, worker3items, maxBucketDelta)
	itsConsistent(its, oldWorker4items, worker4items, maxBucketDelta)
}

func Test_ConsistentHash_notFound_removeBucket(t *testing.T) {
	its := assert.New(t)

	const bucketCount = 5
	const itemCount = 100
	const itemsPerBucket = itemCount / bucketCount
	const maxBucketDelta = (itemsPerBucket / bucketCount) + 4

	var buckets []string
	for x := 0; x < bucketCount; x++ {
		buckets = append(buckets, fmt.Sprintf("worker-%d", x))
	}
	var items []string
	for x := 0; x < itemCount; x++ {
		items = append(items, fmt.Sprintf("google_%d", x))
	}

	ch := New(
		OptBuckets(
			buckets...,
		),
	)

	oldAssignments := ch.Assignments(items...)
	its.False(ch.RemoveBucket("not-worker-0"))
	newAssignments := ch.Assignments(items...)

	assignmentsAreEqual(its, oldAssignments, newAssignments)
}

func Test_ConsistentHash_String(t *testing.T) {
	its := assert.New(t)

	const bucketCount = 5

	var buckets []string
	for x := 0; x < bucketCount; x++ {
		buckets = append(buckets, fmt.Sprintf("worker-%d", x))
	}

	ch := New(
		OptBuckets(
			buckets...,
		),
	)
	its.NotEmpty(ch.String())
}

func Test_ConsistentHash_MarshalJSON(t *testing.T) {
	its := assert.New(t)

	const bucketCount = 5

	var buckets []string
	for x := 0; x < bucketCount; x++ {
		buckets = append(buckets, fmt.Sprintf("worker-%d", x))
	}

	ch := New(
		OptBuckets(
			buckets...,
		),
	)

	output, err := json.Marshal(ch)
	its.Nil(err)
	its.NotEmpty(output)

	var verify []HashedBucket
	err = json.Unmarshal(output, &verify)
	its.Nil(err)
	its.Equal(ch.hashring, verify)
}
