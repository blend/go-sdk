/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"errors"
)

// Is is a helper function that returns if an error is an ex.
//
// It will handle if the err is an exception, a multi-error or a regular error.
// "Isness" is evaluated by if the class of the exception matches the class of the cause
func Is(err interface{}, cause error) bool {
	if err == nil || cause == nil {
		return false
	}
	// handle the case of multi-exceptions first.
	if multiTyped, ok := err.(Multi); ok {
		for _, multiErr := range multiTyped {
			if Is(multiErr, cause) {
				return true
			}
		}
		return false
	}
	if typed := As(err); typed != nil {
		if typed.Class == nil {
			return false
		}
		if causeTyped := As(cause); causeTyped != nil {
			if causeTyped.Class == nil {
				return false
			}
			return typed.Class == causeTyped.Class
		}
		return (typed.Class == cause) || (typed.Class.Error() == cause.Error()) || errors.Is(typed.Class, cause)
	}
	if typed, ok := err.(ClassProvider); ok {
		return typed.Class() == cause || (typed.Class().Error() == cause.Error()) || errors.Is(typed.Class(), cause)
	}

	// handle regular errors
	if typed, ok := err.(error); ok && typed != nil {
		return (err == cause) || (typed.Error() == cause.Error()) || errors.Is(typed, cause)
	}
	// handle ???
	return err == cause
}
