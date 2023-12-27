/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

// CheckValue returns if the value is valid.
func CheckValue(value string) error {
	if len(value) > MaxLabelValueLen {
		return ErrLabelValueTooLong
	}
	return CheckName(value)
}
