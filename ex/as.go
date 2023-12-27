/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

import (
	"errors"
)

// As is a helper method that returns an error as an ex.
func As(err interface{}) *Ex {
	if typed, typedOk := err.(error); typedOk {
		var exx *Ex
		if errors.As(typed, &exx) {
			return exx
		}
		return nil
	}
	if typed, typedOk := err.(Ex); typedOk {
		return &typed
	}
	if typed, typedOk := err.(*Ex); typedOk {
		return typed
	}
	return nil
}
