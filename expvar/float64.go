/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package expvar

import (
	"math"
	"strconv"
	"sync/atomic"
)

// Assert that `Float64` implements `Var`.
var (
	_ Var = (*Float64)(nil)
)

// Float64 is a 64-bit float variable that satisfies the Var interface.
type Float64 struct {
	f uint64
}

// Value returns the underlying float64 value.
func (v *Float64) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&v.f))
}

// String satisfies `Var`.
func (v *Float64) String() string {
	return strconv.FormatFloat(
		math.Float64frombits(atomic.LoadUint64(&v.f)), 'g', -1, 64)
}

// Add adds delta to v.
func (v *Float64) Add(delta float64) {
	for {
		cur := atomic.LoadUint64(&v.f)
		curVal := math.Float64frombits(cur)
		nxtVal := curVal + delta
		nxt := math.Float64bits(nxtVal)
		if atomic.CompareAndSwapUint64(&v.f, cur, nxt) {
			return
		}
	}
}

// Set sets v to value.
func (v *Float64) Set(value float64) {
	atomic.StoreUint64(&v.f, math.Float64bits(value))
}
