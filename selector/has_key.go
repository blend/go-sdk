/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

// HasKey returns if a label set has a given key.
type HasKey string

// Matches returns the selector result.
func (hk HasKey) Matches(labels Labels) bool {
	_, hasKey := labels[string(hk)]
	return hasKey
}

// Validate validates the selector.
func (hk HasKey) Validate() (err error) {
	err = CheckKey(string(hk))
	return
}

// String returns a string representation of the selector.
func (hk HasKey) String() string {
	return string(hk)
}
