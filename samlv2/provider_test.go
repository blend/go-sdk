/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package samlv2_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/samlv2"
)

const signURLPrefix = "https://blend-oie.oktapreview.com/app/blend-oie_samltest_1/exk3vd6k1xkRwxIly1d7/sso/saml?SAMLRequest="

func NewSAMLProvider(audience string) (*samlv2.SAMLProvider, error) {
	metadataRaw, err := os.ReadFile("testdata/metadata.xml")
	if err != nil {
		return nil, err
	}
	config := &samlv2.SAMLConfig{
		AssertionConsumerServiceURL: "http://localhost:8080/saml?redirect_uri=localhost:8081/saml",
		AudienceURI:                 "Audience",
		Metadata:                    string(metadataRaw),
	}

	provider, err := samlv2.New(
		samlv2.OptConfig(config),
		samlv2.OptSkipSignatureValidation(true),
		samlv2.OptValidateEncryptionCert(false),
	)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func Test_SamlResponse(t *testing.T) {
	its := assert.New(t)

	provider, err := NewSAMLProvider("Audience")
	its.Nil(err)

	samlResponse, err := os.ReadFile("testdata/saml_valid.response")
	its.Nil(err)

	assertionInfo, err := provider.OnSAMLResponse(string(samlResponse))
	its.Nil(err)

	its.False(assertionInfo.WarningInfo.NotInAudience)
	its.False(assertionInfo.WarningInfo.InvalidTime)
}

func Test_SamlInvalidTime(t *testing.T) {
	its := assert.New(t)

	provider, err := NewSAMLProvider("Audience")
	its.Nil(err)

	samlResponse, err := os.ReadFile("testdata/saml_invalid.response")
	its.Nil(err)

	assertionInfo, err := provider.OnSAMLResponse(string(samlResponse))
	its.Nil(err)

	its.True(assertionInfo.WarningInfo.InvalidTime)
}

func Test_SamlInvalidAudience(t *testing.T) {
	its := assert.New(t)

	provider, err := NewSAMLProvider("WrongAudience")
	its.Nil(err)

	samlResponse, err := os.ReadFile("testdata/saml_invalid.response")
	its.Nil(err)

	assertionInfo, err := provider.OnSAMLResponse(string(samlResponse))
	its.Nil(err)

	its.True(assertionInfo.WarningInfo.NotInAudience)
}

func Test_BuildURL(t *testing.T) {
	its := assert.New(t)

	provider, err := NewSAMLProvider("WrongAudience")
	its.Nil(err)

	url, err := provider.BuildURL("")
	its.Nil(err)

	its.Equal(strings.HasPrefix(url, signURLPrefix), true)

}

func Test_DefaultCanonicalizer(t *testing.T) {
	its := assert.New(t)

	provider, err := NewSAMLProvider("Audience")
	its.Nil(err)
	its.Equal(provider.Provider.SigningContext().Canonicalizer.Algorithm().String(), "http://www.w3.org/2006/12/xml-c14n11")
}

func Test_ExclusiveCanonicalizer(t *testing.T) {
	its := assert.New(t)
	metadataRaw, err := os.ReadFile("testdata/metadata.xml")
	its.Nil(err)
	config := &samlv2.SAMLConfig{
		AssertionConsumerServiceURL: "http://localhost:8080/saml?redirect_uri=localhost:8081/saml",
		AudienceURI:                 "Audience",
		Metadata:                    string(metadataRaw),
		SigningXMLCanonicalizer:     samlv2.CanonicalXML10ExclusiveAlgorithmID,
	}

	provider, err := samlv2.New(
		samlv2.OptConfig(config),
		samlv2.OptSkipSignatureValidation(true),
		samlv2.OptValidateEncryptionCert(false),
	)
	its.Nil(err)
	its.Equal(provider.Provider.SigningContext().Canonicalizer.Algorithm().String(), "http://www.w3.org/2001/10/xml-exc-c14n#")
}

func Test_InclusiveCanonicalizer(t *testing.T) {
	its := assert.New(t)
	metadataRaw, err := os.ReadFile("testdata/metadata.xml")
	its.Nil(err)
	config := &samlv2.SAMLConfig{
		AssertionConsumerServiceURL: "http://localhost:8080/saml?redirect_uri=localhost:8081/saml",
		AudienceURI:                 "Audience",
		Metadata:                    string(metadataRaw),
		SigningXMLCanonicalizer:     samlv2.CanonicalXML11AlgorithmID,
	}

	provider, err := samlv2.New(
		samlv2.OptConfig(config),
		samlv2.OptSkipSignatureValidation(true),
		samlv2.OptValidateEncryptionCert(false),
	)
	its.Nil(err)
	its.Equal(provider.Provider.SigningContext().Canonicalizer.Algorithm().String(), "http://www.w3.org/2006/12/xml-c14n11")
}

func Test_UnsupportCanonicalizer(t *testing.T) {
	its := assert.New(t)
	metadataRaw, err := os.ReadFile("testdata/metadata.xml")
	its.Nil(err)
	config := &samlv2.SAMLConfig{
		AssertionConsumerServiceURL: "http://localhost:8080/saml?redirect_uri=localhost:8081/saml",
		AudienceURI:                 "Audience",
		Metadata:                    string(metadataRaw),
		SigningXMLCanonicalizer:     "Unsupported Canonicalizer",
	}

	provider, err := samlv2.New(
		samlv2.OptConfig(config),
		samlv2.OptSkipSignatureValidation(true),
		samlv2.OptValidateEncryptionCert(false),
	)
	its.Nil(provider)
	its.True(errors.Is(err, ex.Class(samlv2.ErrorUnsupportedCanonicalizer)))
}
