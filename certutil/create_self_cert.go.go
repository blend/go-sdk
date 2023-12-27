/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

// CreateSelfServerCert creates a self signed server certificate bundle.
func CreateSelfServerCert(commonName string, options ...CertOption) (*CertBundle, error) {
	ca, err := CreateCertificateAuthority()
	if err != nil {
		return nil, err
	}
	return CreateServer(commonName, ca, options...)
}
