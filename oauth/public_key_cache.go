package oauth

import (
	"context"
	"crypto/rsa"
	"net/http"
	"sync"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/jwt"
	"github.com/blend/go-sdk/r2"
)

// PublicKeyCache holds cached signing certs.
type PublicKeyCache struct {
	FetchPublicKeysDefaults []r2.Option
	mu                      sync.RWMutex
	current                 *PublicKeysResponse
}

// Keyfunc returns a jwt keyfunc for a specific exchange tied to context.
func (pkc *PublicKeyCache) Keyfunc(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if token == nil {
			return nil, Error("invalid jwt; token is unset")
		}
		kid, ok := token.Header["kid"]
		if !ok {
			return nil, Error("invalid jwt header; `kid` missing")
		}
		typedKid, ok := kid.(string)
		if !ok {
			return nil, Error("invalid jwt header; `kid` not a string")
		}
		return pkc.Get(ctx, typedKid)
	}
}

// Get gets a cert by id.
func (pkc *PublicKeyCache) Get(ctx context.Context, id string) (*rsa.PublicKey, error) {
	var jwk jwt.JWK
	var ok bool
	pkc.mu.RLock()
	if pkc.current != nil && !pkc.current.IsExpired() {
		jwk, ok = pkc.current.Keys[id]
	}
	pkc.mu.RUnlock()
	if ok {
		return jwk.PublicKey()
	}

	pkc.mu.Lock()
	defer pkc.mu.Unlock()

	// check again after grabbing the lock if
	// the keys have been updated
	if pkc.current != nil && !pkc.current.IsExpired() {
		jwk, ok = pkc.current.Keys[id]
	}
	if ok {
		return jwk.PublicKey()
	}

	// if we should still refresh after grabbing
	// the write lock
	keys, err := FetchPublicKeys(ctx, pkc.FetchPublicKeysDefaults...)
	if err != nil {
		return nil, err
	}
	pkc.current = keys

	jwk, ok = pkc.current.Keys[id]
	if !ok {
		return nil, ex.New("invalid jwt key id; not found in signing keys cache", ex.OptMessagef("Key ID: %s", id))
	}
	return jwk.PublicKey()
}

// PublicKeysResponse is a response for the google certs api.
type PublicKeysResponse struct {
	CacheControl string
	Expires      time.Time
	Keys         map[string]jwt.JWK
}

// IsExpired returns if the cert response is expired.
func (pkr PublicKeysResponse) IsExpired() bool {
	if pkr.Expires.IsZero() {
		return true
	}
	return time.Now().UTC().After(pkr.Expires.UTC())
}

// FetchPublicKeys gets the google signing certs.
func FetchPublicKeys(ctx context.Context, opts ...r2.Option) (*PublicKeysResponse, error) {
	var jwks fetchPublicKeysResponse
	meta, err := r2.New(GoogleKeysURL, opts...).JSON(&jwks)
	if err != nil {
		return nil, err
	}

	expiresHeader := meta.Header.Get(http.CanonicalHeaderKey("Expires"))
	if expiresHeader == "" {
		return nil, ex.New("invalid google keys response; expires unset")
	}

	expires, err := time.Parse(http.TimeFormat, expiresHeader)
	if err != nil {
		return nil, ex.New("invalid google keys response; invalid expires value", ex.OptInner(err))
	}
	res := &PublicKeysResponse{
		Keys:         jwkLookup(jwks.Keys),
		CacheControl: meta.Header.Get(http.CanonicalHeaderKey("Cache-Control")),
		Expires:      expires,
	}
	return res, nil
}

type fetchPublicKeysResponse struct {
	Keys []jwt.JWK `json:"keys"`
}

// jwkLookup creates a jwk lookup.
func jwkLookup(jwks []jwt.JWK) map[string]jwt.JWK {
	output := make(map[string]jwt.JWK)
	for _, jwk := range jwks {
		// We don't check that `jwk.KID` collides with an existing key. We trust that
		// the public certs URL (e.g. the one from Google) does not include duplicates.
		output[jwk.KID] = jwk
	}
	return output
}
