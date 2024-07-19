/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

import (
	"strings"
)

// Strings is a type alias for []string with some helper methods.
// Deprecated: Use collections/generic version.
type Strings []string

// Reverse reverses the strings array in place.
// Deprecated: Use collections/generic version.
func (sa Strings) Reverse() (output Strings) {
	saLen := len(sa)

	switch saLen {
	case 0:
		return
	case 1:
		output = Strings{sa[0]}
		return
	}

	output = make(Strings, len(sa))
	saLen2 := saLen >> 1
	var nx int
	for x := 0; x < saLen2; x++ {
		nx = saLen - (x + 1)
		output[x] = sa[nx]
		output[nx] = sa[x]
	}
	if saLen%2 != 0 {
		output[saLen2] = sa[saLen2]
	}
	return
}

// First returns the first element of the array.
// Deprecated: Use collections/generic version.
func (sa Strings) First() string {
	if len(sa) == 0 {
		return ""
	}
	return sa[0]
}

// Last returns the last element of the array.
// Deprecated: Use collections/generic version.
func (sa Strings) Last() string {
	if len(sa) == 0 {
		return ""
	}
	return sa[len(sa)-1]
}

// Contains returns if the given string is in the array.
// Deprecated: Use collections/generic version.
func (sa Strings) Contains(elem string) bool {
	for _, arrayElem := range sa {
		if arrayElem == elem {
			return true
		}
	}
	return false
}

// ContainsLower returns true if the `elem` is in the Strings, false otherwise.
// Deprecated: Use collections/generic version.
func (sa Strings) ContainsLower(elem string) bool {
	for _, arrayElem := range sa {
		if strings.ToLower(arrayElem) == elem {
			return true
		}
	}
	return false
}

// GetByLower returns an element from the array that matches the input.
// Deprecated: Use collections/generic version.
func (sa Strings) GetByLower(elem string) string {
	for _, arrayElem := range sa {
		if strings.ToLower(arrayElem) == elem {
			return arrayElem
		}
	}
	return ""
}
