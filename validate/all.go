/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

// All returns a validator that returns all non-nil validation errors from a given
// set of validators.
func All(validators ...Validator) Validator {
	return func() error {
		var output []error
		var err error
		for _, validator := range validators {
			if err = validator(); err != nil {
				if errs, hasMany := err.(ValidationErrors); hasMany {
					output = append(output, errs...)
				} else {
					output = append(output, err)
				}
			}
		}
		if len(output) > 0 {
			return ValidationErrors(output)
		}
		return nil
	}
}
