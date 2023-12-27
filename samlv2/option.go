/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package samlv2

import dsig "github.com/russellhaering/goxmldsig"

// OptConfig sets the SAML Provider config.
func OptConfig(cfg *SAMLConfig) Option {
	return func(sp *SAMLProvider) error {
		sp.Config = cfg
		return nil
	}
}

// OptSkipSignatureValidation skips SAML response vaidation.
func OptSkipSignatureValidation(validation bool) Option {
	return func(sp *SAMLProvider) error {
		sp.SkipSignatureValidation = validation
		return nil
	}
}

// OptValidateEncryptionCert sets validatoin of the ecnryption certificate.
func OptValidateEncryptionCert(validate bool) Option {
	return func(sp *SAMLProvider) error {
		sp.ValidateEncryptionCert = validate
		return nil
	}
}

// OptClientKeyStore is used for signing client AuthN requests
func OptClientKeyStore(store dsig.X509KeyStore) Option {
	return func(sp *SAMLProvider) error {
		sp.ClientKeyStore = store
		return nil
	}
}

// Option mutates a SAML Provider.
type Option func(*SAMLProvider) error
