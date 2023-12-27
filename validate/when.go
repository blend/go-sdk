/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

// When returns the result of the "passes" validator if the predicate returns true,
// otherwise it returns nil.
func When(predicate func() bool, passes Validator) Validator {
	return func() error {
		if predicate() {
			return passes()
		}
		return nil
	}
}
