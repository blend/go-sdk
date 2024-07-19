/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import "fmt"

// Equals returns if a key strictly equals a value.
type Equals struct {
	Key, Value string
}

// Matches returns the selector result.
func (e Equals) Matches(labels Labels) bool {
	if value, hasValue := labels[e.Key]; hasValue {
		return e.Value == value
	}
	return false
}

// Validate validates the selector.
func (e Equals) Validate() (err error) {
	err = CheckKey(e.Key)
	if err != nil {
		return
	}
	err = CheckValue(e.Value)
	return
}

// String returns the string representation of the selector.
func (e Equals) String() string {
	return fmt.Sprintf("%s == %s", e.Key, e.Value)
}
