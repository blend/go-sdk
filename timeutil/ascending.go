/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import (
	"sort"
	"time"
)

var (
	_ sort.Interface = (*Ascending)(nil)
)

// Ascending sorts a given list of times ascending, or min to max.
type Ascending []time.Time

// Len implements sort.Sorter
func (a Ascending) Len() int { return len(a) }

// Swap implements sort.Sorter
func (a Ascending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements sort.Sorter
func (a Ascending) Less(i, j int) bool { return a[i].Before(a[j]) }
