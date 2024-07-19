/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
)

// HMAC512 sha512 hashes data with the given key.
func HMAC512(key, plainText []byte) []byte {
	mac := hmac.New(sha512.New, key)
	_, _ = mac.Write([]byte(plainText))
	return mac.Sum(nil)
}

// HMAC256 sha256 hashes data with the given key.
func HMAC256(key, plainText []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(plainText))
	return mac.Sum(nil)
}
