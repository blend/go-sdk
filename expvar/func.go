/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package expvar

import "encoding/json"

// Assert that `Func` implements `Var`.
var (
	_ Var = (*Func)(nil)
)

// Func implements Var by calling the function
// and formatting the returned value using JSON.
type Func func() interface{}

// Value yields the result of calling the function.
func (f Func) Value() interface{} {
	return f()
}

// String implements `Var`.
func (f Func) String() string {
	v, _ := json.Marshal(f())
	return string(v)
}
