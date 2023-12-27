/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"crypto/x509"
	"time"
)

// DefaultOptionsCertificateAuthority are the default options for certificate authorities.
var DefaultOptionsCertificateAuthority = CertOptions{
	Certificate: x509.Certificate{
		IsCA:			true,
		KeyUsage:		x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid:	true,
	},
	NotAfterProvider:	func() time.Time { return time.Now().UTC().AddDate(DefaultCANotAfterYears, 0, 0) },
}

// DefaultOptionsServer are the default create cert options for server certificates.
var DefaultOptionsServer = CertOptions{
	Certificate: x509.Certificate{
		ExtKeyUsage:	[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:	x509.KeyUsageDigitalSignature,
	},
	NotAfterProvider:	func() time.Time { return time.Now().UTC().AddDate(DefaultServerNotAfterYears, 0, 0) },
}

// DefaultOptionsClient are the default create cert options for client certificates.
var DefaultOptionsClient = CertOptions{
	Certificate: x509.Certificate{
		ExtKeyUsage:	[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:	x509.KeyUsageDigitalSignature,
	},
	NotAfterProvider:	func() time.Time { return time.Now().UTC().AddDate(DefaultClientNotAfterYears, 0, 0) },
}
