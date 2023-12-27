/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

// SchemeIsSecure returns if a given scheme is secure.
//
// This is typically used for the `Secure` flag on cookies.
func SchemeIsSecure(scheme string) bool {
	switch scheme {
	case SchemeHTTPS, SchemeSPDY:
		return true
	}
	return false
}
