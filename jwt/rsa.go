package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

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
	hasher.Write([]byte(signingString))

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
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, hasher.Sum(nil))
	if err == nil {
		return EncodeSegment(sigBytes), nil
	}
	return "", err
}

//ParseRSAPrivateKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 private key.
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ex.New(ErrKeyMustBePEMEncoded)
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, ex.New(err)
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ex.New(ErrNotRSAPrivateKey)
	}

	return pkey, nil
}

// ParseRSAPrivateKeyFromPEMWithPassword parses a PEM encoded PKCS1 or PKCS8 private key protected with password.
func ParseRSAPrivateKeyFromPEMWithPassword(key []byte, password string) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ex.New(ErrKeyMustBePEMEncoded)
	}

	var parsedKey interface{}

	var blockDecrypted []byte
	if blockDecrypted, err = x509.DecryptPEMBlock(block, []byte(password)); err != nil {
		return nil, ex.New(err)
	}

	if parsedKey, err = x509.ParsePKCS1PrivateKey(blockDecrypted); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(blockDecrypted); err != nil {
			return nil, ex.New(err)
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ex.New(ErrNotRSAPrivateKey)
	}

	return pkey, nil
}

// ParseRSAPublicKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 public key.
func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ex.New(ErrKeyMustBePEMEncoded)
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, ex.New(err)
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ex.New(ErrNotRSAPublicKey)
	}

	return pkey, nil
}
