package envoyutil

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

const (
	// ErrMissingXFCC is the error returned when XFCC is missing
	ErrMissingXFCC = ex.Class("Missing X-Forwarded-Client-Cert header")
	// ErrInvalidXFCC is the error returned when XFCC is invalid
	ErrInvalidXFCC = ex.Class("Invalid X-Forwarded-Client-Cert header")
	// ErrMissingExtractFunction is the message used when the "extract client
	// identity" function is `nil` or not provided.
	ErrMissingExtractFunction = ex.Class("Missing client identity extraction function")
	// ErrVerifierNil is the message prefix used when a provided verifier is `nil`.
	ErrVerifierNil = ex.Class("XFCC verifier must not be `nil`")
)

// ClientIdentityProvider is a function to extra the client identity from a
// parsed XFCC header. For example, client identity could be determined from the
// SPIFFE URI in the `URI` field in an XFCC element.
type ClientIdentityProvider func(xfcc XFCCElement) (clientIdentity string, err error)

// VerifyXFCC is an "extra" verifier for an XFCC, for example if the server
// identity (from the `By` field in an XFCC element) should be verified in
// addition to the client identity.
type VerifyXFCC func(xfcc XFCCElement) (err *XFCCValidationError)

// ExtractAndVerifyClientIdentity enables extracting client identity from a request.
// It does so by requiring the XFCC header to present and valid and contain exactly
// one element. Then it passes the parsed XFCC header along to some verifiers (e.g.
// to verify the server identity) as well as to an extractor (for the client identity).
func ExtractAndVerifyClientIdentity(req *http.Request, cip ClientIdentityProvider, verifiers ...VerifyXFCC) (string, error) {
	if cip == nil {
		return "", &XFCCFatalError{Class: ErrMissingExtractFunction}
	}

	// Early exit if XFCC header is not present.
	xfccValue := req.Header.Get(HeaderXFCC)
	if xfccValue == "" {
		return "", &XFCCExtractionError{Class: ErrMissingXFCC}
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
