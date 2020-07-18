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
//       - `extractJustURI` satisfies `envoyutil.ClientIdentityProvider`.
//       - `extractFailure` satisfies `envoyutil.ClientIdentityProvider`.
var (
	_ envoyutil.ClientIdentityProvider = extractJustURI
	_ envoyutil.ClientIdentityProvider = extractFailure
)

func TestExtractAndVerifyClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		XFCC           string
		ClientIdentity string
		ErrorType      string
		Class          ex.Class
		Extract        envoyutil.ClientIdentityProvider
		Verifiers      []envoyutil.VerifyXFCC
	}
	testCases := []testCase{
		{ErrorType: "XFCCFatalError", Class: envoyutil.ErrMissingExtractFunction},
		{XFCC: "", ErrorType: "XFCCExtractionError", Class: envoyutil.ErrMissingXFCC, Extract: extractJustURI},
		{XFCC: `""`, ErrorType: "XFCCValidationError", Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
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

func TestClientIdentityFromSPIFFE(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		XFCC           string
		TrustDomain    string
		ClientIdentity string
		ErrorType      string
		Class          ex.Class
		Denied         []string
	}
	testCases := []testCase{
		{
			XFCC:      "URI=not-spiffe",
			ErrorType: "XFCCExtractionError",
			Class:     envoyutil.ErrInvalidClientIdentity,
		},
		{
			XFCC:           "URI=spiffe://cluster.local/ns/light1/sa/bulb",
			TrustDomain:    "cluster.local",
			ClientIdentity: "bulb.light1",
		},
		{
			XFCC:        "URI=spiffe://cluster.local/ns/light2/sa/bulb",
			TrustDomain: "k8s.local",
			ErrorType:   "XFCCValidationError",
			Class:       envoyutil.ErrInvalidClientIdentity,
		},
		{
			XFCC:        "URI=spiffe://cluster.local/ns/light3/sa/bulb/extra",
			TrustDomain: "cluster.local",
			ErrorType:   "XFCCExtractionError",
			Class:       envoyutil.ErrInvalidClientIdentity,
		},
		{
			XFCC:        "URI=spiffe://cluster.local/ns/light4/sa/bulb",
			TrustDomain: "cluster.local",
			ErrorType:   "XFCCValidationError",
			Class:       envoyutil.ErrDeniedClientIdentity,
			Denied:      []string{"bulb.light4"},
		},
		{
			XFCC:           "URI=spiffe://cluster.local/ns/light5/sa/bulb",
			TrustDomain:    "cluster.local",
			ClientIdentity: "bulb.light5",
			Denied:         []string{"not.me"},
		},
		{
			XFCC:        "URI=spiffe://cluster.local/ns/light6/sa/bulb",
			TrustDomain: "cluster.local",
			ErrorType:   "XFCCValidationError",
			Class:       envoyutil.ErrDeniedClientIdentity,
			Denied:      []string{"not.me", "bulb.light6", "also.not-me"},
		},
	}

	for _, tc := range testCases {
		xfccElements, err := envoyutil.ParseXFCC(tc.XFCC)
		assert.Nil(err)
		assert.Len(xfccElements, 1)
		xfcc := xfccElements[0]

		cip := envoyutil.ClientIdentityFromSPIFFE(tc.TrustDomain, tc.Denied...)
		clientIdentity, err := cip(xfcc)
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
		default:
			assert.Nil(err, tc)
		}
	}
}

// extractJustURI satisfies `envoyutil.ClientIdentityProvider` and just returns the URI.
func extractJustURI(xfcc envoyutil.XFCCElement) (string, error) {
	return xfcc.URI, nil
}

// extractFailure satisfies `envoyutil.ClientIdentityProvider` and fails.
func extractFailure(xfcc envoyutil.XFCCElement) (string, error) {
	return "", &envoyutil.XFCCExtractionError{Class: "extractFailure", XFCC: xfcc.String()}
}

func makeVerifyXFCC(expectedBy string) envoyutil.VerifyXFCC {
	return func(xfcc envoyutil.XFCCElement) *envoyutil.XFCCValidationError {
		if xfcc.By == expectedBy {
			return nil
		}

		c := ex.Class(fmt.Sprintf("verifyFailure: expected %q", expectedBy))
		return &envoyutil.XFCCValidationError{Class: c, XFCC: xfcc.String()}
	}
}
