/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

// ReturnFirst runs a list of validators and returns
// the first validator to error (if there is one).
func ReturnFirst(validators ...Validator) error {
	var err error
	for _, validator := range validators {
		if err = validator(); err != nil {
			return err
		}
	}
	return nil
}
