package envoyutil_test

import (
	"fmt"
	"net/http"
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"

	"github.com/blend/go-sdk/envoyutil"
)

// NOTE: Ensure
//       - `extractJustURI` satisfies `envoyutil.ExtractFromXFCC`.
//       - `extractFailure` satisfies `envoyutil.ExtractFromXFCC`.
var (
	_ envoyutil.ExtractFromXFCC = extractJustURI
	_ envoyutil.ExtractFromXFCC = extractFailure
)

func TestExtractClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		XFCC           string
		ClientIdentity string
		ErrorType      string
		Class          ex.Class
		Extract        envoyutil.ExtractFromXFCC
		Verifiers      []envoyutil.VerifyXFCC
	}
	testCases := []testCase{
		{ErrorType: "XFCCFatalError", Class: envoyutil.MissingExtractFunction},
		{XFCC: "", ErrorType: "XFCCExtractionError", Class: envoyutil.ErrMissingXFCC, Extract: extractJustURI},
		{XFCC: `""`, ErrorType: "XFCCExtractionError", Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
		{XFCC: "something=bad", ErrorType: "XFCCExtractionError", Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
		{
			XFCC:           "By=spiffe://cluster.local/ns/blend/sa/idea;URI=spiffe://cluster.local/ns/light/sa/bulb",
			ClientIdentity: "spiffe://cluster.local/ns/light/sa/bulb",
			Extract:        extractJustURI,
		},
		{XFCC: "By=x;URI=y", ErrorType: "XFCCExtractionError", Class: "extractFailure", Extract: extractFailure},
		{
			XFCC:      "By=abc;URI=def",
			ErrorType: "XFCCExtractionError",
			Class:     `verifyFailure: expected "xyz"`,
			Extract:   extractJustURI,
			Verifiers: []envoyutil.VerifyXFCC{makeVerifyXFCC("xyz")},
		},
		{
			XFCC:           "By=abc;URI=def",
			ClientIdentity: "def",
			Extract:        extractJustURI,
			Verifiers:      []envoyutil.VerifyXFCC{makeVerifyXFCC("abc")},
		},
		{
			XFCC:      "By=abc;URI=def",
			ErrorType: "XFCCFatalError",
			Class:     envoyutil.VerifierNil,
			Extract:   extractJustURI,
			Verifiers: []envoyutil.VerifyXFCC{nil},
		},
	}

	for _, tc := range testCases {
		// Set-up mock context.
		r, newReqErr := http.NewRequest("GET", "", nil)
		assert.Nil(newReqErr)
		if tc.XFCC != "" {
			r.Header.Add(envoyutil.HeaderXFCC, tc.XFCC)
		}

		clientIdentity, err := envoyutil.ExtractClientIdentity(r, tc.Extract, tc.Verifiers...)
		assert.Equal(tc.ClientIdentity, clientIdentity)
		switch tc.ErrorType {
		case "XFCCExtractionError":
			expected := &envoyutil.XFCCExtractionError{Class: tc.Class, XFCC: tc.XFCC}
			assert.Equal(expected, err)
		case "XFCCFatalError":
			expected := &envoyutil.XFCCFatalError{Class: tc.Class, XFCC: tc.XFCC}
			assert.Equal(expected, err)
		default:
			assert.Nil(err)
		}
	}
}

// extractJustURI satisfies `envoyutil.ExtractFromXFCC` and just returns the URI.
func extractJustURI(xfcc envoyutil.XFCCElement, _ string) (string, *envoyutil.XFCCExtractionError) {
	return xfcc.URI, nil
}

// extractFailure satisfies `envoyutil.ExtractFromXFCC` and fails.
func extractFailure(xfcc envoyutil.XFCCElement, xfccValue string) (string, *envoyutil.XFCCExtractionError) {
	return "", &envoyutil.XFCCExtractionError{Class: "extractFailure", XFCC: xfccValue}
}

func makeVerifyXFCC(expectedBy string) envoyutil.VerifyXFCC {
	return func(xfcc envoyutil.XFCCElement, xfccValue string) *envoyutil.XFCCValidationError {
		if xfcc.By == expectedBy {
			return nil
		}

		c := ex.Class(fmt.Sprintf("verifyFailure: expected %q", expectedBy))
		return &envoyutil.XFCCValidationError{Class: c, XFCC: xfccValue}
	}
}
