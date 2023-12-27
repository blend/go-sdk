/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
	Expired	RemovalReason	= iota
	Removed	RemovalReason	= iota
)
