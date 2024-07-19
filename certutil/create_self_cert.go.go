/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
