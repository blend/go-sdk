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
//       - `extractJustURI` satisfies `envoyutil.IdentityProvider`.
//       - `extractFailure` satisfies `envoyutil.IdentityProvider`.
var (
	_ envoyutil.IdentityProvider = extractJustURI
	_ envoyutil.IdentityProvider = extractFailure
)

func TestExtractAndVerifyClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		XFCC           string
		ClientIdentity string
		ErrorType      string
		Class          ex.Class
		Extract        envoyutil.IdentityProvider
		Verifiers      []envoyutil.VerifyXFCC
	}
	testCases := []testCase{
		{ErrorType: "XFCCFatalError", Class: envoyutil.ErrMissingExtractFunction},
		{XFCC: "", ErrorType: "XFCCExtractionError", Class: envoyutil.ErrMissingXFCC, Extract: extractJustURI},
		{XFCC: `""`, ErrorType: "XFCCExtractionError", Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
		{XFCC: "something=bad", ErrorType: "XFCCExtractionError", Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
		{
			XFCC:      "By=first,URI=second",
			ErrorType: "XFCCValidationError",
			Class:     envoyutil.ErrInvalidXFCC,
			Extract:   extractJustURI,
		},
		{
			XFCC:           "By=spiffe://cluster.local/ns/blend/sa/idea;URI=spiffe://cluster.local/ns/light/sa/bulb",
			ClientIdentity: "spiffe://cluster.local/ns/light/sa/bulb",
			Extract:        extractJustURI,
		},
		{XFCC: "By=x;URI=y", ErrorType: "XFCCExtractionError", Class: "extractFailure", Extract: extractFailure},
		{
			XFCC:      "By=abc;URI=def",
			ErrorType: "XFCCValidationError",
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
			Class:     envoyutil.ErrVerifierNil,
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

		clientIdentity, err := envoyutil.ExtractAndVerifyClientIdentity(r, tc.Extract, tc.Verifiers...)
		assert.Equal(tc.ClientIdentity, clientIdentity)
		switch tc.ErrorType {
		case "XFCCExtractionError":
			assert.True(envoyutil.IsExtractionError(err), tc)
			expected := &envoyutil.XFCCExtractionError{Class: tc.Class, XFCC: tc.XFCC}
			assert.Equal(expected, err, tc)
		case "XFCCValidationError":
			assert.True(envoyutil.IsValidationError(err), tc)
			expected := &envoyutil.XFCCValidationError{Class: tc.Class, XFCC: tc.XFCC}
			assert.Equal(expected, err, tc)
		case "XFCCFatalError":
			assert.True(envoyutil.IsFatalError(err), tc)
			expected := &envoyutil.XFCCFatalError{Class: tc.Class, XFCC: tc.XFCC}
			assert.Equal(expected, err, tc)
		default:
			assert.Nil(err, tc)
		}
	}
}

// extractJustURI satisfies `envoyutil.IdentityProvider` and just returns the URI.
func extractJustURI(xfcc envoyutil.XFCCElement) (string, error) {
	return xfcc.URI, nil
}

// extractFailure satisfies `envoyutil.IdentityProvider` and fails.
func extractFailure(xfcc envoyutil.XFCCElement) (string, error) {
	return "", &envoyutil.XFCCExtractionError{Class: "extractFailure", XFCC: xfcc.String()}
}

func makeVerifyXFCC(expectedBy string) envoyutil.VerifyXFCC {
	return func(xfcc envoyutil.XFCCElement) error {
		if xfcc.By == expectedBy {
			return nil
		}

		c := ex.Class(fmt.Sprintf("verifyFailure: expected %q", expectedBy))
		return &envoyutil.XFCCValidationError{Class: c, XFCC: xfcc.String()}
	}
}
