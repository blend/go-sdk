/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"errors"
)

// Is is a helper function that returns if an error is an ex.
//
// It will handle if the err is an exception, a multi-error or a regular error.
// "Isness" is evaluated by if the class of the exception matches the class of the cause.
//
// Deprecated: Use [errors.Is]. Make sure `Is()` and `Unwrap()` are properly
// implemented on your custom classes.
func Is(err error, cause error) bool {
	if err == nil || cause == nil {
		return false
	}
	// If it's a ClassProvider, try comparing with the class first.
	if typed, ok := err.(ClassProvider); ok && Is(typed.Class(), cause) {
		return true
	}
	// Otherwise, use the native `errors.Is()`. Ex and Multi errors are handled in
	// this case as they implement the `Is` method.
	return errors.Is(err, cause)
}
