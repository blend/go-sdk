/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

// Reverse returns a new slice with the elements reversed.
// Reverse returns nil if the slice is empty.
func Reverse[T any](s []T) []T {
	var res []T
	for i := len(s) - 1; i >= 0; i-- {
		res = append(res, s[i])
	}
	return res
}

// ReverseInPlace reverses the generic slice in place.
func ReverseInPlace[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// First returns the first element of the slice.
func First[T any](s []T) T {
	var res T
	if len(s) == 0 {
		return res
	}
	return s[0]
}

// Last returns the last element of the slice.
func Last[T any](s []T) T {
	var res T
	if len(s) == 0 {
		return res
	}
	return s[len(s)-1]
}

// Contains returns if the given string is in the slice.
func Contains[T comparable](s []T, target T) bool {
	for _, e := range s {
		if e == target {
			return true
		}
	}
	return false
}
