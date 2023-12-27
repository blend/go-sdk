/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrSomeNil		ex.Class	= "all references should be nil"
	ErrSomeNotNil		ex.Class	= "at least one reference should not be nil"
	ErrSomeOneNotNil	ex.Class	= "exactly one of a set of reference should not be nil"
)

// Many defines validator rules that apply to a variadic
// set of untyped references.
func Many(objs ...interface{}) ManyValidators {
	return ManyValidators{Objs: objs}
}

// ManyValidators returns the validator singleton for some rules.
type ManyValidators struct {
	Objs []interface{}
}

// Nil ensures none of a set of references aren't nil.
func (mv ManyValidators) Nil() Validator {
	return func() error {
		for _, obj := range mv.Objs {
			if !IsNil(obj) {
				return Error(ErrSomeNil, nil)
			}
		}
		return nil
	}
}

// NotNil ensures at least one reference is set.
func (mv ManyValidators) NotNil() Validator {
	return func() error {
		for _, obj := range mv.Objs {
			if !IsNil(obj) {
				return nil
			}
		}
		return Error(ErrSomeNotNil, nil)
	}
}

// OneNotNil ensures  at least and at most one reference is set.
func (mv ManyValidators) OneNotNil() Validator {
	return func() error {
		var set bool
		for _, obj := range mv.Objs {
			if !IsNil(obj) {
				if set {
					return Error(ErrSomeOneNotNil, nil)
				}
				set = true
			}
		}
		if !set {
			return Error(ErrSomeOneNotNil, nil)
		}
		return nil
	}
}
