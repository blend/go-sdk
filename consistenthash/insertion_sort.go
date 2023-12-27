/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package consistenthash

import "sort"

// InsertionSort inserts an bucket into a hashring by binary searching
// for the index which would satisfy the overall "sorted" status of the ring
// returning the updated hashring.
func InsertionSort(ring []HashedBucket, item HashedBucket) []HashedBucket {
	destination := sort.Search(len(ring), func(index int) bool {
		return ring[index].Hashcode >= item.Hashcode
	})
	ring = append(ring, HashedBucket{})
	copy(ring[destination+1:], ring[destination:])
	ring[destination] = item
	return ring
}
