/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package consistenthash

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_insertionSort(t *testing.T) {
	its := assert.New(t)

	ring := []HashedBucket{}
	ring = InsertionSort(ring, HashedBucket{Hashcode: 3})
	ring = InsertionSort(ring, HashedBucket{Hashcode: 1})
	ring = InsertionSort(ring, HashedBucket{Hashcode: 4})
	ring = InsertionSort(ring, HashedBucket{Hashcode: 2})
	ring = InsertionSort(ring, HashedBucket{Hashcode: 0})
	ring = InsertionSort(ring, HashedBucket{Hashcode: 5})

	its.Len(ring, 6)

	its.Equal(0, ring[0].Hashcode)
	its.Equal(1, ring[1].Hashcode)
	its.Equal(2, ring[2].Hashcode)
	its.Equal(3, ring[3].Hashcode)
	its.Equal(4, ring[4].Hashcode)
	its.Equal(5, ring[5].Hashcode)
}
