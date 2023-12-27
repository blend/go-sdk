/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

// OptAfter appends after run actions.
func OptAfter(steps ...SuiteAction) Option {
	return func(s *Suite) {
		s.After = append(s.After, steps...)
	}
}
