/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

// First is a validator that returns the first error of a given set of validators.
func First(validators ...Validator) Validator {
	return func() error {
		var err error
		for _, validator := range validators {
			if err = validator(); err != nil {
				return err
			}
		}
		return nil
	}
}
