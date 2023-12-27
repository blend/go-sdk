/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
