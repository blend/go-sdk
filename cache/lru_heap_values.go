/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

import (
	"container/heap"
)

var (
	_ heap.Interface = (*LRUHeapValues)(nil)
)

// LRUHeapValues is an alias that allows for use of a Value array as a heap
// storage backend. It also implements sorting for other use cases.
type LRUHeapValues []*Value

// Len returns the values length.
func (lruv LRUHeapValues) Len() int	{ return len(lruv) }

// Less returns if two values are strictly less than eachother.
func (lruv LRUHeapValues) Less(i, j int) bool	{ return lruv[i].Expires.Before(lruv[j].Expires) }

// Swap swaps values at the given positions.
func (lruv LRUHeapValues) Swap(i, j int)	{ lruv[i], lruv[j] = lruv[j], lruv[i] }

// Push adds a new item.
func (lruv *LRUHeapValues) Push(x interface{}) {
	*lruv = append(*lruv, x.(*Value))
}

// Pop returns an item.
func (lruv *LRUHeapValues) Pop() interface{} {
	old := *lruv
	n := len(old)
	x := old[n-1]
	*lruv = old[0 : n-1]
	return x
}
