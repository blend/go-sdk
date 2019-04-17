package crypto

import (
	cryptorand "crypto/rand"
	"encoding/hex"

	"github.com/blend/go-sdk/ex"
)

// MustCreateKey creates a key, if an error is returned, it panics.
func MustCreateKey(keySize int) []byte {
	key, err := CreateKey(keySize)
	if err != nil {
		panic(err)
	}
	return key
}

// CreateKey creates a key of a given size by reading that much data off the crypto/rand reader.
func CreateKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := cryptorand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// MustCreateKeyString generates a new key and returns it as a string.
func MustCreateKeyString() string {
	return hex.EncodeToString(MustCreateKey(KeySize))
}

// ParseKey parses a key from a string.
func ParseKey(key string) ([]byte, error) {
	decoded, err := hex.DecodeString(key)
	if err != nil {
		return nil, ex.New(err)
	}
	if len(decoded) != KeySize {
		return nil, ex.New("parse key; invalid key length")
	}
	return decoded, nil
}
