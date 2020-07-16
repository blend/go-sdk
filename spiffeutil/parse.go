package spiffeutil

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

const (
	spiffePrefix = "spiffe://"

	// ErrInvalidURI is the class of error returned when parsing SPIFFE URI fails
	ErrInvalidURI = ex.Class("Invalid SPIFFE URI")
)

// Parse consumes a SPIFFE URI and splits out the trust domain and workload
// identifier. For example in `spiffe://cluster.local/ns/blend/sa/quasar`
// the trust domain is `cluster.local` and the workload identifier is
// `ns/blend/sa/quasar`.
func Parse(uri string) (string, string, error) {
	if !strings.HasPrefix(uri, spiffePrefix) {
		return "", "", ex.New(ErrInvalidURI).WithMessagef("Does not match protocol: %q", uri)
	}

	suffix := uri[len(spiffePrefix):]
	parts := strings.SplitN(suffix, "/", 2)
	if len(parts) != 2 || len(parts[1]) == 0 {
		return "", "", ex.New(ErrInvalidURI).WithMessagef("Missing workload identifier: %q", uri)
	}

	return parts[0], parts[1], nil
}
