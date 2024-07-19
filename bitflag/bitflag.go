/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bitflag

// Combine combines all the values into one flag.
func Combine(values ...Bitflag) Bitflag {
	var outputFlag uint64
	for _, value := range values {
		outputFlag = outputFlag | uint64(value)
	}
	return Bitflag(outputFlag)
}

// Bitflag is a large unsighted integer flag set.
type Bitflag uint64

// All returns if all of a given set of flags are set.
func (bf Bitflag) All(values Bitflag) bool {
	return uint64(bf)&uint64(values) == uint64(values)
}

// Any returns if any the reference bits are set for a given value
func (bf Bitflag) Any(values Bitflag) bool {
	return uint64(bf)&uint64(values) > 0
}

// Set sets a flag value to 1.
func (bf Bitflag) Set(values Bitflag) Bitflag {
	return Bitflag(uint64(bf) | uint64(values))
}

// Unset makes a given flag zero'd in the set.
func (bf Bitflag) Unset(values Bitflag) Bitflag {
	return Bitflag(uint64(bf) ^ ((-(0) ^ uint64(values)) & uint64(bf)))
}
