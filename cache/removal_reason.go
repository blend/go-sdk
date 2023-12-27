/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cache

// RemovalReason is a reason for removal.
type RemovalReason int

// String returns a string representation of the removal reason.
func (rr RemovalReason) String() string {
	switch rr {
	case Expired:
		return "expired"
	case Removed:
		return "removed"
	default:
		return "unknown"
	}
}

// RemovalReasons
const (
	Expired RemovalReason = iota
	Removed RemovalReason = iota
)
