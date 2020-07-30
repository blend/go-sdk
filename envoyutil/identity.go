package envoyutil

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

const (
	// ErrMissingXFCC is the error returned when XFCC is missing.
	ErrMissingXFCC = ex.Class("Missing X-Forwarded-Client-Cert header")
	// ErrInvalidXFCC is the error returned when XFCC is invalid.
	ErrInvalidXFCC = ex.Class("Invalid X-Forwarded-Client-Cert header")
	// ErrInvalidClientIdentity is the error returned when XFCC has a
	// missing / invalid client identity.
	ErrInvalidClientIdentity = ex.Class("Client identity could not be determined from X-Forwarded-Client-Cert header")
	// ErrDeniedClientIdentity is the error returned when a parsed client identity is in a deny list or
	// not in an allow list.
	ErrDeniedClientIdentity = ex.Class("Client identity from X-Forwarded-Client-Cert header is denied")
	// ErrInvalidServerIdentity is the error returned when XFCC has a
	// missing / invalid client identity.
	ErrInvalidServerIdentity = ex.Class("Server identity could not be determined from X-Forwarded-Client-Cert header")
	// ErrDeniedServerIdentity is the error returned when a parsed client identity is in a deny list or
	// not in an allow list.
	ErrDeniedServerIdentity = ex.Class("Server identity from X-Forwarded-Client-Cert header is denied")
	// ErrMissingExtractFunction is the message used when the "extract client
	// identity" function is `nil` or not provided.
	ErrMissingExtractFunction = ex.Class("Missing client identity extraction function")
	// ErrVerifierNil is the message prefix used when a provided verifier is `nil`.
	ErrVerifierNil = ex.Class("XFCC verifier must not be `nil`")
)

// IdentityProvider is a function to extract the client or server identity from
// a parsed XFCC header. For example, client identity could be determined from the
// SPIFFE URI in the `URI` field in an XFCC element.
type IdentityProvider func(xfcc XFCCElement) (identity string, err error)

// VerifyXFCC is an "extra" verifier for an XFCC, for example if the server
// identity (from the `By` field in an XFCC element) should be verified in
// addition to the client identity.
type VerifyXFCC func(xfcc XFCCElement) error

// ExtractAndVerifyClientIdentity enables extracting client identity from a request.
// It does so by requiring the XFCC header to be present and valid and contain exactly
// one element. Then it passes the parsed XFCC header along to some `verifiers` (e.g.
// to verify the server identity) as well as to an extractor `cip` (for the
// client identity).
func ExtractAndVerifyClientIdentity(req *http.Request, cip IdentityProvider, verifiers ...VerifyXFCC) (string, error) {
	if cip == nil {
		return "", &XFCCFatalError{Class: ErrMissingExtractFunction}
	}

	// Early exit if XFCC header is not present.
	xfccValue := req.Header.Get(HeaderXFCC)
	if xfccValue == "" {
		return "", &XFCCValidationError{Class: ErrMissingXFCC}
	}

	// Early exit if XFCC header is invalid, or has zero or multiple elements.
	xfccElements, parseErr := ParseXFCC(xfccValue)
	if parseErr != nil {
		return "", &XFCCExtractionError{Class: ErrInvalidXFCC, XFCC: xfccValue}
	}
	if len(xfccElements) != 1 {
		return "", &XFCCValidationError{Class: ErrInvalidXFCC, XFCC: xfccValue}
	}
	xfcc := xfccElements[0]

	// Run all verifiers on the parsed `xfcc`.
	for _, verifier := range verifiers {
		if verifier == nil {
			return "", &XFCCFatalError{Class: ErrVerifierNil, XFCC: xfccValue}
		}

		err := verifier(xfcc)
		if err != nil {
			return "", err
		}
	}

	// Do final extraction.
	return cip(xfcc)
}
