/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package jwk

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

type (
	// Set represents a set of JWKs as defined by https://tools.ietf.org/html/rfc7517#section-5
	Set	struct {
		Keys []JWK `json:"keys"`
	}

	// JWK represents a cryptographic key as defined by https://tools.ietf.org/html/rfc7517#section-4
	JWK	struct {
		KTY	string	`json:"kty"`
		USE	string	`json:"use,omitempty"`
		ALG	string	`json:"alg,omitempty"`
		KID	string	`json:"kid,omitempty"`
		E	string	`json:"e,omitempty"`
		N	string	`json:"n,omitempty"`
	}
)

// RSAPublicKey parses the public key in the JWK to a rsa.PublicKey.
func (j JWK) RSAPublicKey() (*rsa.PublicKey, error) {
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
		E:	int(e.Int64()),
		N:	&n,
	}, nil
}

// KTY parameter values as defined in https://tools.ietf.org/html/rfc7518#section-6.1
const (
	KTYRSA = "RSA"
)

// RSAPublicKeyToJWK converts an RSA public key to a JWK.
func RSAPublicKeyToJWK(key *rsa.PublicKey) JWK {
	return JWK{
		KTY:	KTYRSA,
		E:	base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
		N:	base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
	}
}
