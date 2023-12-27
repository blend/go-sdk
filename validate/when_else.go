/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

// WhenElse returns the result of the "passes" validator if the predicate returns true,
// otherwise it returns the result of the "fails" validator.
func WhenElse(predicate func() bool, passes, fails Validator) Validator {
	return func() error {
		if predicate() {
			if passes != nil {
				return passes()
			}
			return nil
		}
		if fails != nil {
			return fails()
		}
		return nil
	}
}
