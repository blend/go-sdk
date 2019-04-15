package jwt

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Parser is a parser for tokens.
type Parser struct {
	ValidMethods         []string // If populated, only these methods will be considered valid
	UseJSONNumber        bool     // Use JSON Number format in JSON decoder
	SkipClaimsValidation bool     // Skip claims validation during token parsing
}

// Parse parses, validate, and return a token.
func (p *Parser) Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
	return p.ParseWithClaims(tokenString, MapClaims{}, keyFunc)
}

// ParseWithClaims parses a token with a given set of claims.
func (p *Parser) ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {
	token, parts, err := p.ParseUnverified(tokenString, claims)
	if err != nil {
		return token, err
	}

	// Verify signing method is in the required set
	if p.ValidMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range p.ValidMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			return token, ex.New(ErrValidation, ex.OptInner(ErrInvalidSigningMethod))
		}
	}

	// Lookup key
	var key interface{}
	if keyFunc == nil {
		// keyFunc was not provided.  short circuiting validation
		return token, ex.New(ErrValidation, ex.OptInner(ErrKeyfuncUnset))
	}

	if key, err = keyFunc(token); err != nil {
		return token, err
	}

	// Validate Claims
	if !p.SkipClaimsValidation {
		if err := token.Claims.Valid(); err != nil {
			// this is strictly an aud, exp, or nbf style validation error.
			return token, ex.New(ErrValidation, ex.OptInner(err))
		}
	}

	// Perform validation
	token.Signature = parts[2]
	if err = token.Method.Verify(strings.Join(parts[0:2], "."), token.Signature, key); err != nil {
		return token, ex.New(ErrValidation, ex.OptInner(ex.New(ErrValidationSignature, ex.OptInner(err))))
	}

	token.Valid = true
	return token, nil
}

// ParseUnverified parses the token but doesn't validate the signature.
// WARNING: Don't use this method unless you know what you're doing
// It's only ever useful in cases where you know the signature is valid
// (because it has been checked previously in the stack) and you want to extract values from it.
func (p *Parser) ParseUnverified(tokenString string, claims Claims) (token *Token, parts []string, err error) {
	parts = strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, parts, ex.New(ErrValidation, ex.OptMessage("token contains an invalid number of segments"))
	}

	token = &Token{Raw: tokenString}

	// parse Header
	var headerBytes []byte
	if headerBytes, err = DecodeSegment(parts[0]); err != nil {
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			return token, parts, ex.New(ErrValidation, ex.OptMessage("tokenstring should not contain 'bearer '"))
		}
		return token, parts, ex.New(ErrValidation, ex.OptInner(err))
	}
	if err = json.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, parts, ex.New(ErrValidation, ex.OptInner(err))
	}

	// parse Claims
	var claimBytes []byte
	token.Claims = claims

	if claimBytes, err = DecodeSegment(parts[1]); err != nil {
		return token, parts, ex.New(ErrValidation, ex.OptInner(err))
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	if p.UseJSONNumber {
		dec.UseNumber()
	}
	// JSON Decode.  Special case for map type to avoid weird pointer behavior
	if c, ok := token.Claims.(MapClaims); ok {
		err = dec.Decode(&c)
	} else {
		err = dec.Decode(&claims)
	}
	// Handle decode error
	if err != nil {
		return token, parts, ex.New(ErrValidation, ex.OptInner(err))
	}

	// Lookup signature method
	if method, ok := token.Header["alg"].(string); ok {
		if token.Method = GetSigningMethod(method); token.Method == nil {
			return token, parts, ex.New(ErrValidation, ex.OptInner(ErrInvalidSigningMethod))
		}
		return token, parts, nil
	}
	return token, parts, ex.New(ErrValidation, ex.OptInner(ErrInvalidSigningMethod))
}
