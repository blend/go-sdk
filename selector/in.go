/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"fmt"
	"strings"
)

// In returns if a key matches a set of values.
type In struct {
	Key    string
	Values []string
}

// Matches returns the selector result.
func (i In) Matches(labels Labels) bool {
	// if the labels has a given key
	if value, hasValue := labels[i.Key]; hasValue {
		// for each selector value
		for _, iv := range i.Values {
			// if they match, return true
			if iv == value {
				return true
			}
		}
		return false
	}
	// in should be exclusive, that is
	// we should fail if the in key isn't present
	return false
}

// Validate validates the selector.
func (i In) Validate() (err error) {
	err = CheckKey(i.Key)
	if err != nil {
		return
	}
	for _, v := range i.Values {
		err = CheckValue(v)
		if err != nil {
			return
		}
	}
	return
}

// String returns a string representation of the selector.
func (i In) String() string {
	return fmt.Sprintf("%s in (%s)", i.Key, strings.Join(i.Values, ", "))
}
