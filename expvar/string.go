/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package expvar

import (
	"encoding/json"
	"sync/atomic"
)

// Assert that `String` implements `Var`.
var (
	_ Var = (*String)(nil)
)

// String is a string variable, and satisfies the Var interface.
type String struct {
	s atomic.Value	// string
}

// Value returns the underlying value.
func (v *String) Value() string {
	p, _ := v.s.Load().(string)
	return p
}

// String implements the Var interface. To get the unquoted string
// use Value.
func (v *String) String() string {
	s := v.Value()
	b, _ := json.Marshal(s)
	return string(b)
}

// Set sets the value
func (v *String) Set(value string) {
	v.s.Store(value)
}
