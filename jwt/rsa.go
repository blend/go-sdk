package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/blend/go-sdk/ex"
)

// SigningMethodRSA implements the RSA family of signing methods signing methods
// Expects *rsa.PrivateKey for signing and *rsa.PublicKey for validation
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

// Alg returns the name of the signing method.
func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

// Verify implements the Verify method from SigningMethod
// For this signing method, must be an *rsa.PublicKey structure.
func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = DecodeSegment(signature); err != nil {
		return err
	}

	var rsaKey *rsa.PublicKey
	var ok bool

	if rsaKey, ok = key.(*rsa.PublicKey); !ok {
		return ex.New(ErrInvalidKeyType)
	}

	// Create hasher
	if !m.Hash.Available() {
		return ex.New(ErrHashUnavailable)
	}
	hasher := m.Hash.New()
	if _, err := hasher.Write([]byte(signingString)); err != nil {
		return ex.New(err)
	}

	// Verify the signature
	return ex.New(rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig))
}

// Sign implements the Sign method from SigningMethod
// For this signing method, must be an *rsa.PrivateKey structure.
func (m *SigningMethodRSA) Sign(signingString string, key interface{}) (string, error) {
	var rsaKey *rsa.PrivateKey
	var ok bool

	// Validate type of key
	if rsaKey, ok = key.(*rsa.PrivateKey); !ok {
		return "", ErrInvalidKey
	}

	// Create the hasher
	if !m.Hash.Available() {
		return "", ex.New(ErrHashUnavailable)
	}

	hasher := m.Hash.New()
	if _, err := hasher.Write([]byte(signingString)); err != nil {
		return "", ex.New(err)
	}

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, hasher.Sum(nil))
	if err != nil {
		return "", ex.New(err)
	}
	return EncodeSegment(sigBytes), nil
}
