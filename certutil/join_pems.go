/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import "strings"

// JoinPEMs appends pem blocks together with newlines.
//
// Each pem block will have `strings.TrimSpace()` called on it.
//
// Usage note: you should add pems in the following order:
// - leaf
// - intermediate
// - root
// It's a little baffling, basically the other way around from what you'd thing probably.
func JoinPEMs(pems ...string) string {
	var cleaned []string
	for _, pem := range pems {
		pemCleaned := strings.TrimSpace(pem)
		if pemCleaned != "" {
			cleaned = append(cleaned, pemCleaned)
		}
	}
	return strings.Join(cleaned, "\n") + "\n"
}
