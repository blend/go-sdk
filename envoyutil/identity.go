package envoyutil

import (
	"net/http"
)

const (
	// ErrMissingExtractFunction is the error returned when the "extract client
	// identity" function is `nil` or not provided.
	ErrMissingExtractFunction = "Missing client identity extraction function"
	// ErrVerifierNil is the error returned when a provided verifier is `nil`.
	ErrVerifierNil = "Verifier must not be `nil`"
	// ErrMissingXFCC is the error returned when XFCC is missing
	ErrMissingXFCC = "Missing X-Forwarded-Client-Cert header"
	// ErrInvalidXFCC is the error returned when XFCC is invalid
	ErrInvalidXFCC = "Invalid X-Forwarded-Client-Cert header"
)

// InvalidXFCCResponse contains metadata about an XFCC header that could not
// be parsed. This is intended to be used as the body of a 401 Unauthorized
// response.
type InvalidXFCCResponse struct {
	Message  string      `json:"message" xml:"message"`
	XFCC     string      `json:"xfcc,omitempty" xml:"xfcc,omitempty"`
	Metadata interface{} `json:"metadata,omitempty" xml:"metadata,omitempty"`
}

// ExtractFromXFCC is a function to extra the client identity from a
// parsed XFCC header. For example, client identity could be determined from the
// SPIFFE URI in the `URI` field in an XFCC element.
type ExtractFromXFCC func(xfcc XFCCElement, xfccValue string) (clientIdentity string, errResponse *InvalidXFCCResponse)

// VerifyXFCC is an "extra" verifier for an XFCC, for example if the server
// identity (from the `By` field in an XFCC element) should be verified in
// addition to the client identity.
type VerifyXFCC func(xfcc XFCCElement, xfccValue string) (errResponse *InvalidXFCCResponse)

// ExtractClientIdentity enables extracting client identity from a request. It
// does so by requiring the XFCC header to present and valid and then passing
// the parsed XFCC header along to some verifiers (e.g. to verify the server
// identity) as well as to an extractor (for the client identity).
func ExtractClientIdentity(req *http.Request, efx ExtractFromXFCC, verifiers ...VerifyXFCC) (clientIdentity string, errResponse *InvalidXFCCResponse) {
	if efx == nil {
		errResponse = &InvalidXFCCResponse{Message: ErrMissingExtractFunction}
		return
	}

	// Early exit if XFCC header is not present.
	xfccValue := req.Header.Get(HeaderXFCC)
	if xfccValue == "" {
		errResponse = &InvalidXFCCResponse{Message: ErrMissingXFCC}
		return
	}

	// Early exit if XFCC header is invalid, or has zero or multiple elements.
	xfccElements, err := ParseXFCC(xfccValue)
	if err != nil || len(xfccElements) != 1 {
		errResponse = &InvalidXFCCResponse{Message: ErrInvalidXFCC, XFCC: xfccValue}
		return
	}
	xfcc := xfccElements[0]

	// Run all verifiers on the parsed `xfcc`.
	for _, verifier := range verifiers {
		if verifier == nil {
			errResponse = &InvalidXFCCResponse{Message: ErrVerifierNil, XFCC: xfccValue}
			return
		}

		errResponse = verifier(xfcc, xfccValue)
		if errResponse != nil {
			return
		}
	}

	// Do final extraction.
	return efx(xfcc, xfccValue)
}
