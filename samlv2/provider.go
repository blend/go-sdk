/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package samlv2

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

// X509KeyStore is a store keeping references to public/private signing keys
type X509KeyStore struct {
	privateKey *rsa.PrivateKey
	cert       []byte
}

// GetKeyPair returns public/private key pair from a store
func (ks *X509KeyStore) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	return ks.privateKey, ks.cert, nil
}

// SAMLProvider is wrapper dedicated for
// verification and validation of SAML assertion documents.
type SAMLProvider struct {
	//Config references SAML configuration
	Config *SAMLConfig
	//Log is the default logger
	Log logger.Logger
	//SkipSignatureValidation skips validating SAML response signature
	SkipSignatureValidation bool
	//ValidateEncryptionCert validates signature certificates if set to true
	ValidateEncryptionCert bool
	//Provider is SAMLv2 service provider
	Provider *saml2.SAMLServiceProvider
	//ClientKeyStore to sign Authn requests
	ClientKeyStore dsig.X509KeyStore
}

// ParseMetadata parses SAML IDP metadata, extracts basic SAML attributes
// and certificates for SAML provider
func ParseMetadata(config *SAMLConfig) (*x509.Certificate, error) {
	metadata := &types.EntityDescriptor{}
	err := xml.Unmarshal([]byte(config.Metadata), metadata)
	if err != nil {
		return nil, err
	}

	var idpCert *x509.Certificate
	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				return nil, fmt.Errorf("metadata certificate(%d) must not be empty", idx)
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				return nil, err
			}

			idpCert, err = x509.ParseCertificate(certData)
			if err != nil {
				return nil, err
			}
		}
	}

	config.IdentityProviderSSOURL = metadata.IDPSSODescriptor.SingleSignOnServices[0].Location
	config.IdentityProviderIssuer = metadata.EntityID
	config.ServiceProviderIssuer = metadata.EntityID

	return idpCert, nil
}

// New returns a new SAML provider.
func New(opts ...Option) (*SAMLProvider, error) {
	// creates SAMLV2 service provider with default signature validation
	// set to true
	p := &SAMLProvider{
		ValidateEncryptionCert:  true,
		SkipSignatureValidation: false,
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, ex.New(err)
		}
	}

	if p.Config == nil {
		return nil, ex.New("SAML config should be provided")
	}

	if p.Config.Metadata == "" {
		return nil, ex.New("identity provider metadata is a required parameter")
	}

	certStore := &dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	idpCert, err := ParseMetadata(p.Config)
	if err != nil {
		return nil, ex.New(err)
	}

	certStore.Roots = append(certStore.Roots, idpCert)

	// We sign the AuthnRequest with a random key because Okta doesn't seem
	// to verify these.
	clientKeyStore := dsig.RandomKeyStoreForTest()

	if p.Config.PublicKey != nil {
		// Parse PEM block
		var block *pem.Block
		if block, _ = pem.Decode(p.Config.PrivateKey); block == nil {
			return nil, errors.New("must have no nil private key")
		}

		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		if block, _ = pem.Decode(p.Config.PublicKey); block == nil {
			return nil, errors.New("must have no nil public key")
		}

		pubCertBytes := block.Bytes
		pubCert, err := x509.ParseCertificate(pubCertBytes)
		if err != nil {
			return nil, err
		}

		// we're completely replacing certificate Roots here
		certStore.Roots = []*x509.Certificate{pubCert}

		//Note: we overwrite client key store in a config here
		//in case public/private key pair is defined in a config
		p.ClientKeyStore = &X509KeyStore{
			privateKey: privKey,
			cert:       pubCertBytes,
		}
	}

	if p.ClientKeyStore != nil {
		clientKeyStore = p.ClientKeyStore
	}
	var xmlCanonicalizer dsig.Canonicalizer
	switch p.Config.SigningXMLCanonicalizer {
	case CanonicalXML10ExclusiveAlgorithmID:
		xmlCanonicalizer = dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("")
	case CanonicalXML11AlgorithmID, "":
		xmlCanonicalizer = dsig.MakeC14N11Canonicalizer()
	default:
		return nil, ex.New(ex.Class(ErrorUnsupportedCanonicalizer))
	}

	p.Provider = &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:         p.Config.IdentityProviderSSOURL,
		IdentityProviderIssuer:         p.Config.IdentityProviderIssuer,
		ServiceProviderIssuer:          p.Config.IdentityProviderIssuer,
		AssertionConsumerServiceURL:    p.Config.AssertionConsumerServiceURL,
		SignAuthnRequests:              true,
		SignAuthnRequestsCanonicalizer: xmlCanonicalizer,
		AudienceURI:                    p.Config.AudienceURI,
		IDPCertificateStore:            certStore,
		SPKeyStore:                     clientKeyStore,
		ValidateEncryptionCert:         p.ValidateEncryptionCert,
		SkipSignatureValidation:        p.SkipSignatureValidation,
	}

	return p, nil
}

// BuildURL creates SAML Auth URL
func (p *SAMLProvider) BuildURL(state string) (string, error) {
	return p.Provider.BuildAuthURL(state)
}

// OnSAMLResponse decodes, validates and verifies SAML Assertion Response
func (p *SAMLProvider) OnSAMLResponse(response string) (*saml2.AssertionInfo, error) {
	assertionInfo, err := p.Provider.RetrieveAssertionInfo(response)
	if err != nil {
		return nil, ex.New(err)
	}

	return assertionInfo, nil
}
