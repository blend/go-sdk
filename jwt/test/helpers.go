package test

import (
	"crypto/rsa"
	"github.com/blend/go-sdk/jwt"
)

// MustLoadRSAPrivateKey loads an rsa private key and panics on error.
func MustLoadRSAPrivateKey(keyData []byte) *rsa.PrivateKey {
	key, e := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

// MustLoadRSAPublicKey loads an rsa public key and panics on error.
func MustLoadRSAPublicKey(keyData []byte) *rsa.PublicKey {
	key, e := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

// MakeSampleToken makes a sample token.
func MakeSampleToken(c jwt.Claims, key interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	s, e := token.SignedString(key)

	if e != nil {
		panic(e.Error())
	}

	return s
}
