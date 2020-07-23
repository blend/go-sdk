package crypto

import (
	"crypto/hmac"
	"crypto/sha512"
)

// HMAC512 sha512 hashes data with the given key.
func HMAC512(key, plainText []byte) []byte {
	mac := hmac.New(sha512.New, key)
	_, _ = mac.Write([]byte(plainText))
	return mac.Sum(nil)
}
