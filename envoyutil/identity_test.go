package envoyutil_test

import (
	"fmt"
	"net/http"
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"

	"github.com/blend/go-sdk/envoyutil"
)

// NOTE: Ensure
//       - `extractJustURI` satisfies `envoyutil.ExtractFromXFCC`.
//       - `extractFailure` satisfies `envoyutil.ExtractFromXFCC`.
var (
	_ envoyutil.ExtractFromXFCC = extractJustURI
	_ envoyutil.ExtractFromXFCC = extractFailure
)

func TestXFCCErrorError(t *testing.T) {
	assert := sdkAssert.New(t)

	c := "oh a bad thing happened"
	var err error = &envoyutil.XFCCError{Class: c}
	assert.Equal(c, err.Error())
}

func TestExtractClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		XFCC           string
		ClientIdentity string
		Class          string
		Extract        envoyutil.ExtractFromXFCC
		Verifiers      []envoyutil.VerifyXFCC
	}
	testCases := []testCase{
		{Class: envoyutil.ErrMissingExtractFunction},
		{XFCC: "", Class: envoyutil.ErrMissingXFCC, Extract: extractJustURI},
		{XFCC: `""`, Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
		{XFCC: "something=bad", Class: envoyutil.ErrInvalidXFCC, Extract: extractJustURI},
		{XFCC: "By=spiffe://cluster.local/ns/blend/sa/idea;URI=spiffe://cluster.local/ns/light/sa/bulb", ClientIdentity: "spiffe://cluster.local/ns/light/sa/bulb", Extract: extractJustURI},
		{XFCC: "By=x;URI=y", Class: "extractFailure", Extract: extractFailure},
		{
			XFCC:      "By=abc;URI=def",
			Verifiers: []envoyutil.VerifyXFCC{makeVerifyXFCC("xyz")},
			Class:     `verifyFailure: expected "xyz"`,
			Extract:   extractJustURI,
		},
		{
			XFCC:           "By=abc;URI=def",
			Verifiers:      []envoyutil.VerifyXFCC{makeVerifyXFCC("abc")},
			Extract:        extractJustURI,
			ClientIdentity: "def",
		},
		{
			XFCC:      "By=abc;URI=def",
			Verifiers: []envoyutil.VerifyXFCC{nil},
			Class:     envoyutil.ErrVerifierNil,
			Extract:   extractJustURI,
		},
	}

	for _, tc := range testCases {
		// Set-up mock context.
		r, err := http.NewRequest("GET", "", nil)
		assert.Nil(err)
		if tc.XFCC != "" {
			r.Header.Add(envoyutil.HeaderXFCC, tc.XFCC)
		}

		clientIdentity, errResponse := envoyutil.ExtractClientIdentity(r, tc.Extract, tc.Verifiers...)
		assert.Equal(tc.ClientIdentity, clientIdentity)
		if tc.Class == "" {
			assert.Nil(errResponse)
		} else {
			expected := &envoyutil.XFCCError{Class: tc.Class, XFCC: tc.XFCC}
			assert.Equal(expected, errResponse)
		}
	}
}

// extractJustURI satisfies `envoyutil.ExtractFromXFCC` and just returns the URI.
func extractJustURI(xfcc envoyutil.XFCCElement, _ string) (string, *envoyutil.XFCCError) {
	return xfcc.URI, nil
}

// extractFailure satisfies `envoyutil.ExtractFromXFCC` and fails.
func extractFailure(xfcc envoyutil.XFCCElement, xfccValue string) (string, *envoyutil.XFCCError) {
	return "", &envoyutil.XFCCError{Class: "extractFailure", XFCC: xfccValue}
}

func makeVerifyXFCC(expectedBy string) envoyutil.VerifyXFCC {
	return func(xfcc envoyutil.XFCCElement, xfccValue string) *envoyutil.XFCCError {
		if xfcc.By == expectedBy {
			return nil
		}

		c := fmt.Sprintf("verifyFailure: expected %q", expectedBy)
		return &envoyutil.XFCCError{Class: c, XFCC: xfccValue}
	}
}
