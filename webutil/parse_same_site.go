package webutil

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

// mustNil requires that an error is `nil`.
func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

// MustParseSameSite parses a string value for same site and panics on error.
func MustParseSameSite(sameSite string) http.SameSite {
	value, err := ParseSameSite(sameSite)
	mustNil(err)
	return value
}

// ParseSameSite parses a string value for same site.
func ParseSameSite(sameSite string) (http.SameSite, error) {
	switch sameSite {
	case SameSiteStrict:
		return http.SameSiteStrictMode, nil
	case SameSiteLax:
		return http.SameSiteLaxMode, nil
	case SameSiteDefault:
		return http.SameSiteDefaultMode, nil
	default:
		return http.SameSite(-1), ex.New(ErrInvalidSameSite, ex.OptMessagef("value: %s", sameSite))
	}
}
