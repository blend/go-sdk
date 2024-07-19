/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

// OptAfter appends after run actions.
func OptAfter(steps ...SuiteAction) Option {
	return func(s *Suite) {
		s.After = append(s.After, steps...)
	}
}
