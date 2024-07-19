/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
//   - `extractJustURI` satisfies `envoyutil.IdentityProvider`.
//   - `extractFailure` satisfies `envoyutil.IdentityProvider`.
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
		{XFCC: "", ErrorType: "XFCCValidationError", Class: envoyutil.ErrMissingXFCC, Extract: extractJustURI},
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

func TestSPIFFEClientIdentityProvider(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		XFCC           string
		TrustDomain    string
		ClientIdentity string
		ErrorType      string
		Class          ex.Class
		Metadata       interface{}
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
			Metadata:    map[string]string{"trustDomain": "cluster.local"},
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
			Metadata:    map[string]string{"clientIdentity": "bulb.light4"},
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
			Metadata:    map[string]string{"clientIdentity": "bulb.light6"},
			Denied:      []string{"not.me", "bulb.light6", "also.not-me"},
		},
	}

	for _, tc := range testCases {
		xfccElements, err := envoyutil.ParseXFCC(tc.XFCC)
		assert.Nil(err)
		assert.Len(xfccElements, 1)
		xfcc := xfccElements[0]

		cip := envoyutil.SPIFFEClientIdentityProvider(
			envoyutil.OptAllowedTrustDomains(tc.TrustDomain),
			envoyutil.OptDeniedIdentities(tc.Denied...),
		)
		clientIdentity, err := cip(xfcc)
		assert.Equal(tc.ClientIdentity, clientIdentity)

		switch tc.ErrorType {
		case "XFCCExtractionError":
			assert.True(envoyutil.IsExtractionError(err), tc)
			expected := &envoyutil.XFCCExtractionError{Class: tc.Class, XFCC: tc.XFCC, Metadata: tc.Metadata}
			assert.Equal(expected, err, tc)
		case "XFCCValidationError":
			assert.True(envoyutil.IsValidationError(err), tc)
			expected := &envoyutil.XFCCValidationError{Class: tc.Class, XFCC: tc.XFCC, Metadata: tc.Metadata}
			assert.Equal(expected, err, tc)
		default:
			assert.Nil(err, tc)
		}
	}
}

func TestSPIFFEServerIdentityProvider(t *testing.T) {
	assert := sdkAssert.New(t)

	// Verifier returns `nil` error when server identity is valid.
	verifier := envoyutil.SPIFFEServerIdentityProvider()
	xfcc := envoyutil.XFCCElement{By: "spiffe://cluster.local/ns/time/sa/line"}
	err := verifier(xfcc)
	assert.Nil(err)

	// Verifier returns extraction error when server identity is invalid.
	verifier = envoyutil.SPIFFEServerIdentityProvider()
	xfcc = envoyutil.XFCCElement{By: "not-spiffe"}
	err = verifier(xfcc)
	assert.True(envoyutil.IsExtractionError(err))
	var expected error = &envoyutil.XFCCExtractionError{
		Class: envoyutil.ErrInvalidServerIdentity,
		XFCC:  xfcc.String(),
	}
	assert.Equal(expected, err)

	// Verifier returns validation error when server identity is in deny list.
	verifier = envoyutil.SPIFFEServerIdentityProvider(
		envoyutil.OptDeniedIdentities("line.time"),
	)
	xfcc = envoyutil.XFCCElement{By: "spiffe://cluster.local/ns/time/sa/line"}
	err = verifier(xfcc)
	assert.True(envoyutil.IsValidationError(err))
	expected = &envoyutil.XFCCValidationError{
		Class:    envoyutil.ErrDeniedServerIdentity,
		XFCC:     xfcc.String(),
		Metadata: map[string]string{"serverIdentity": "line.time"},
	}
	assert.Equal(expected, err)
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
