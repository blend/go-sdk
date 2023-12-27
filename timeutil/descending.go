/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package timeutil

import (
	"sort"
	"time"
)

var (
	_ sort.Interface = (*Ascending)(nil)
)

// Descending sorts a given list of times ascending, or min to max.
type Descending []time.Time

// Len implements sort.Sorter
func (d Descending) Len() int	{ return len(d) }

// Swap implements sort.Sorter
func (d Descending) Swap(i, j int)	{ d[i], d[j] = d[j], d[i] }

// Less implements sort.Sorter
func (d Descending) Less(i, j int) bool	{ return d[i].After(d[j]) }
