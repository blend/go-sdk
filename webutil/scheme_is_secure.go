/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
