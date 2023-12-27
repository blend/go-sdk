/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

// OptBefore appends before run actions.
func OptBefore(steps ...SuiteAction) Option {
	return func(s *Suite) {
		s.Before = append(s.Before, steps...)
	}
}
