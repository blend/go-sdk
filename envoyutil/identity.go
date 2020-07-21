package envoyutil

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

const (
	// ErrMissingExtractFunction is the error returned when the "extract client
	// identity" function is `nil` or not provided.
	ErrMissingExtractFunction = ex.Class("Missing client identity extraction function")
	// ErrVerifierNil is the error returned when a provided verifier is `nil`.
	ErrVerifierNil = ex.Class("Verifier must not be `nil`")
	// ErrMissingXFCC is the error returned when XFCC is missing
	ErrMissingXFCC = ex.Class("Missing X-Forwarded-Client-Cert header")
	// ErrInvalidXFCC is the error returned when XFCC is invalid
	ErrInvalidXFCC = ex.Class("Invalid X-Forwarded-Client-Cert header")
)

// NOTE: Ensure
//       - `XFCCError` satisfies `error`.
var (
	_ error = (*XFCCError)(nil)
)

// XFCCError contains metadata about an XFCC header that could not
// be parsed. This is intended to be used as the body of a 401 Unauthorized
// response.
type XFCCError struct {
	// Class can be used to uniquely identify the type of the error.
	Class ex.Class `json:"class" xml:"class"`
	// XFCC contains the XFCC header value that could not be parsed or was
	// invalid in some way.
	XFCC string `json:"xfcc,omitempty" xml:"xfcc,omitempty"`
	// Metadata contains extra information relevant to a specific failure.
	Metadata interface{} `json:"metadata,omitempty" xml:"metadata,omitempty"`
}

// Error satisfies the `error` interface. It is intended to be a unique
// identifier for the error.
func (xe *XFCCError) Error() string {
	return string(xe.Class)
}

// ExtractFromXFCC is a function to extra the client identity from a
// parsed XFCC header. For example, client identity could be determined from the
// SPIFFE URI in the `URI` field in an XFCC element.
type ExtractFromXFCC func(xfcc XFCCElement, xfccValue string) (clientIdentity string, err *XFCCError)

// VerifyXFCC is an "extra" verifier for an XFCC, for example if the server
// identity (from the `By` field in an XFCC element) should be verified in
// addition to the client identity.
type VerifyXFCC func(xfcc XFCCElement, xfccValue string) (err *XFCCError)

// ExtractClientIdentity enables extracting client identity from a request. It
// does so by requiring the XFCC header to present and valid and then passing
// the parsed XFCC header along to some verifiers (e.g. to verify the server
// identity) as well as to an extractor (for the client identity).
func ExtractClientIdentity(req *http.Request, efx ExtractFromXFCC, verifiers ...VerifyXFCC) (string, error) {
	if efx == nil {
		return "", &XFCCError{Class: ErrMissingExtractFunction}
	}

	// Early exit if XFCC header is not present.
	xfccValue := req.Header.Get(HeaderXFCC)
	if xfccValue == "" {
		return "", &XFCCError{Class: ErrMissingXFCC}
	}

	// Early exit if XFCC header is invalid, or has zero or multiple elements.
	xfccElements, parseErr := ParseXFCC(xfccValue)
	if parseErr != nil || len(xfccElements) != 1 {
		return "", &XFCCError{Class: ErrInvalidXFCC, XFCC: xfccValue}
	}
	xfcc := xfccElements[0]

	// Run all verifiers on the parsed `xfcc`.
	for _, verifier := range verifiers {
		if verifier == nil {
			return "", &XFCCError{Class: ErrVerifierNil, XFCC: xfccValue}
		}

		err := verifier(xfcc, xfccValue)
		if err != nil {
			return "", err
		}
	}

	// Do final extraction.
	return efx(xfcc, xfccValue)
}
