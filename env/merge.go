/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

// Merge merges a given set of environment variables.
func Merge(sets ...Vars) Vars {
	output := Vars{}
	for _, set := range sets {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}
