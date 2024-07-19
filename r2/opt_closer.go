/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

// OptCloser sets the request closer.
//
// It is typically used to clean up or trigger other actions.
func OptCloser(action func() error) Option {
	return func(r *Request) error {
		r.Closer = action
		return nil
	}
}
