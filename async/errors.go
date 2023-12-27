/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "github.com/blend/go-sdk/ex"

// Errors is a channel for errors
type Errors chan error

// First returns the first (non-nil) error.
func (e Errors) First() error {
	if errorCount := len(e); errorCount > 0 {
		var err error
		for x := 0; x < errorCount; x++ {
			err = <-e
			if err != nil {
				return ex.New(err)
			}
		}
		return nil
	}
	return nil
}

// All returns all the non-nil errors in the channel
// as a multi-error.
func (e Errors) All() error {
	if errorCount := len(e); errorCount > 0 {
		var errors []error
		for x := 0; x < errorCount; x++ {
			err := <-e
			if err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) > 0 {
			return ex.Append(nil, errors...)
		}
		return nil
	}
	return nil
}
