/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package semver

import "github.com/blend/go-sdk/ex"

// GreaterOrEqualTo returns a validator that enforces input versions
// are greater than or equal to a given version.
//
// Strictly speaking, it returns an error if an input version is
// less than the given version.
func GreaterOrEqualTo(version string) func(string) error {
	compiled, err := NewVersion(version)
	if err != nil {
		panic(err)
	}
	return func(compare string) error {
		compareCompiled, err := NewVersion(compare)
		if err != nil {
			return err
		}
		if compiled.LessThan(compareCompiled) {
			return ex.New(ErrConstraintFailed, ex.OptMessagef("greater than or equal to: %v, compare: %s", version, compare))
		}
		return nil
	}
}
