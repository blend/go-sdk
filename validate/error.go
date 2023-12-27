/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"fmt"

	"github.com/blend/go-sdk/ex"
)

// Error returns a new validation error.
// The root class of the error will be ErrValidation.
// The root stack will begin the frame above this call to error.
// The inner error will the cause of the validation vault.
func Error(cause error, value interface{}, messageArgs ...interface{}) error {
	return &ex.Ex{
		Class:      ErrValidation,
		StackTrace: ex.Callers(ex.DefaultNewStartDepth + 1),
		Inner: &ValidationError{
			Cause:   cause,
			Value:   value,
			Message: fmt.Sprint(messageArgs...),
		},
	}
}
