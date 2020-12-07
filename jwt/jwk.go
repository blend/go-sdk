package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

// JWK is a signing cert response.
type JWK struct {
	KTY string `json:"kty"`
	ALG string `json:"alg"`
	USE string `json:"use"`
	KID string `json:"kid"`
	E   string `json:"e"`
	N   string `json:"n"`
}

// PublicKey parses the public key in the JWK.
func (j JWK) PublicKey() (*rsa.PublicKey, error) {
	decodedE, err := base64.RawURLEncoding.DecodeString(j.E)
	if err != nil {
		return nil, err
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(j.N)
	if err != nil {
		return nil, err
	}

	var n, e big.Int
	e.SetBytes(decodedE)
	n.SetBytes(decodedN)
	return &rsa.PublicKey{
		E: int(e.Int64()),
		N: &n,
	}, nil
}
