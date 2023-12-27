/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package samlv2

// SAMLConfig is the SAML provider configuration.
type SAMLConfig struct {
	// Identity Provider SSO URL
	IdentityProviderSSOURL string

	// Identity Provider Issuer
	IdentityProviderIssuer string

	// Service Provider Issuer
	ServiceProviderIssuer string

	// Metadata is metadata.xml that has all the above
	// attributes plus signing IDP certificate.
	Metadata string

	//PublicKey certificate is an optional Root CA PEM certificate
	//gets added to Saml Provider's list of root CAs
	PublicKey []byte

	//PrivateKey is an optional private PEM certificate
	//counterpart to a PublicKey.
	//Public/PrivateKey is used to re-encrypt SAML response
	PrivateKey []byte

	//Assertion Consumer Service URL
	AssertionConsumerServiceURL string

	// SAML Audience
	AudienceURI string

	// Canonicalization Algorithm for XML signing
	SigningXMLCanonicalizer CanonicalizerAlgorithm
}
