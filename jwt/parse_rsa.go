package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/blend/go-sdk/ex"
)

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
