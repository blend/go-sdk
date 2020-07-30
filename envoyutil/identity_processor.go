package envoyutil

import (
	"github.com/blend/go-sdk/spiffeutil"
)

// IdentityProcessorOption mutates an identity processor.
type IdentityProcessorOption func(*IdentityProcessor)

// IdentityFormatter describes functions that will produce an identity string
// from a parsed SPIFFE URI.
type IdentityFormatter = func(XFCCElement, *spiffeutil.ParsedURI) (string, error)

// IdentityType represents the type of identity that will be extracted by an
// `IdentityProcessor`. It can either be a client or server identity.
type IdentityType int

const (
	// ClientIdentity represents client identity.
	ClientIdentity IdentityType = 0
	// ServerIdentity represents server identity.
	ServerIdentity IdentityType = 1
)

// IdentityProcessor provides configurable fields that can be used to
// help validate a parsed SPIFFE URI and produce and validate an identity from
// a parsed SPIFFE URI. The `Type` field determines if a client or server
// identity should be provided; by default the type will be client identity.
type IdentityProcessor struct {
	Type           IdentityType
	FormatIdentity IdentityFormatter
}
