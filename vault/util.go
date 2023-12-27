/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

// ErrClassForStatus returns the exception class for a given remote status code.
func ErrClassForStatus(statusCode int) ex.Class {
	switch statusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden, http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return ErrServerError
	}
}
