/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

import "fmt"

// NotEquals returns if a key strictly equals a value.
type NotEquals struct {
	Key, Value string
}

// Matches returns the selector result.
func (ne NotEquals) Matches(labels Labels) bool {
	if value, hasValue := labels[ne.Key]; hasValue {
		return ne.Value != value
	}
	return true
}

// Validate validates the selector.
func (ne NotEquals) Validate() (err error) {
	err = CheckKey(ne.Key)
	if err != nil {
		return
	}
	err = CheckValue(ne.Value)
	return
}

// String returns a string representation of the selector.
func (ne NotEquals) String() string {
	return fmt.Sprintf("%s != %s", ne.Key, ne.Value)
}
