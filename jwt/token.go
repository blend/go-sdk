package jwt

import (
	"encoding/json"
	"strings"
)

// New creates a new Token with a given signing method.
func New(method SigningMethod) *Token {
	return NewWithClaims(method, MapClaims{})
}

// NewWithClaims creates a new token with the given claims object.
func NewWithClaims(method SigningMethod, claims Claims) *Token {
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}

// Token is a JWT token.
// Different fields will be used depending on whether you're creating or parsing/verifying a token.
type Token struct {
	Raw       string                 // The raw token.  Populated when you Parse a token
	Method    SigningMethod          // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    Claims                 // The second segment of the token
	Signature string                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

// SignedString returns the complete, signed token.
func (t *Token) SignedString(key interface{}) (output string, err error) {
	var sig, sstr string
	if sstr, err = t.signingString(); err != nil {
		return
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return
	}
	output = strings.Join([]string{sstr, sig}, ".")
	return
}

// SigningString generates the signing string.
func (t *Token) signingString() (output string, err error) {
	parts := make([]string, 2)
	for i := range parts {
		var jsonValue []byte
		if i == 0 {
			if jsonValue, err = json.Marshal(t.Header); err != nil {
				return
			}
		} else {
			if jsonValue, err = json.Marshal(t.Claims); err != nil {
				return
			}
		}

		parts[i] = EncodeSegment(jsonValue)
	}
	output = strings.Join(parts, ".")
	return
}
