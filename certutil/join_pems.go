/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
