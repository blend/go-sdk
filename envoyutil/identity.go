package envoyutil

import (
	"net/http"

	"github.com/blend/go-sdk/ex"
)

const (
	// ErrNilInterface is the error returned when a user error results in `nil`
	// being passed as a function interface type.
	ErrNilInterface = ex.Class("User error: nil interface provided")
	// ErrMissingXFCC is the error returned when XFCC is missing
	ErrMissingXFCC = ex.Class("Missing X-Forwarded-Client-Cert header")
	// ErrInvalidXFCC is the error returned when XFCC is invalid
	ErrInvalidXFCC = ex.Class("Invalid X-Forwarded-Client-Cert header")

	// MissingExtractFunction is the message used when the "extract client
	// identity" function is `nil` or not provided.
	MissingExtractFunction = "Missing client identity extraction function"
	// VerifierNil is the message prefix used when a provided verifier is `nil`.
	VerifierNil = "XFCC verifier must not be `nil`"
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
	return xe.Class.Error()
}

// IsUserError determines if an error is a user error (vs. a legitimate input
// error) in this package. A user error happens when a `nil` interface is
// provided in a place where a non-`nil` value is required.
func IsUserError(err error) bool {
	asEx, ok := err.(*ex.Ex)
	if !ok {
		return false
	}
	return asEx.Class == ErrNilInterface
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
		return "", ex.New(ErrNilInterface, ex.OptMessage(MissingExtractFunction))
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
			return "", ex.New(
				ErrNilInterface,
				ex.OptMessagef("%s; XFCC: %q", VerifierNil, xfccValue),
			)
		}

		err := verifier(xfcc, xfccValue)
		if err != nil {
			return "", err
		}
	}

	// Do final extraction.
	return efx(xfcc, xfccValue)
}
