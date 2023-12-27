/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

// Any matches everything
type Any struct{}

// Matches returns true
func (a Any) Matches(labels Labels) bool {
	return true
}

// Validate validates the selector
func (a Any) Validate() (err error) {
	return nil
}

// String returns a string representation of the selector
func (a Any) String() string {
	return ""
}
