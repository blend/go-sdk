/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

// CheckValue returns if the value is valid.
func CheckValue(value string) error {
	if len(value) > MaxLabelValueLen {
		return ErrLabelValueTooLong
	}
	return CheckName(value)
}
