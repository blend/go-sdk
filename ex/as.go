/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"errors"
)

// As is a helper method that returns an error as an ex.
//
// Deprecated: Use [errors.As] with [*Ex]. Make sure `As()` and `Unwrap()` are
// properly implemented on your custom classes.
func As(err interface{}) *Ex {
	switch typed := err.(type) {
	case error:
		var exx *Ex
		if errors.As(typed, &exx) {
			return exx
		}
		return nil
	case Ex:
		return &typed
	default:
		return nil
	}
}
