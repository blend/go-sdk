/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package expvar

import (
	"strconv"
	"sync/atomic"
)

// Assert that `Func` implements `Var`.
var (
	_ Var = (*Int)(nil)
)

// Int is a 64-bit integer variable that satisfies the Var interface.
type Int struct {
	i int64
}

// Value returns the current value.
func (v *Int) Value() int64 {
	return atomic.LoadInt64(&v.i)
}

// String satisfies `Var`.
func (v *Int) String() string {
	return strconv.FormatInt(atomic.LoadInt64(&v.i), 10)
}

// Add adds to the value.
func (v *Int) Add(delta int64) int64 {
	return atomic.AddInt64(&v.i, delta)
}

// Set sets the value
func (v *Int) Set(value int64) {
	atomic.StoreInt64(&v.i, value)
}
