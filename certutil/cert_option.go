/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"os"
	"time"

	"github.com/blend/go-sdk/ex"
)

// CertOption is an option for creating certs.
type CertOption func(*CertOptions) error

// OptSubjectCommonName sets the subject common name.
func OptSubjectCommonName(commonName string) CertOption {
	return func(csr *CertOptions) error {
		csr.Subject.CommonName = commonName
		return nil
	}
}

// OptSubjectAlternateNames sets the subject alternate names.
func OptSubjectAlternateNames(dnsNames ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.DNSNames = dnsNames
		return nil
	}
}

// OptSubjectOrganization sets the subject organization names.
func OptSubjectOrganization(organization ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.Subject.Organization = organization
		return nil
	}
}

// OptSubjectCountry sets the subject country names.
func OptSubjectCountry(country ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.Subject.Country = country
		return nil
	}
}

// OptSubjectProvince sets the subject province names.
func OptSubjectProvince(province ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.Subject.Province = province
		return nil
	}
}

// OptSubjectLocality sets the subject locality names.
func OptSubjectLocality(locality ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.Subject.Locality = locality
		return nil
	}
}

// OptNotAfter sets the not after time.
func OptNotAfter(notAfter time.Time) CertOption {
	return func(csr *CertOptions) error {
		csr.NotAfter = notAfter
		return nil
	}
}

// OptNotBefore sets the not before time.
func OptNotBefore(notBefore time.Time) CertOption {
	return func(csr *CertOptions) error {
		csr.NotBefore = notBefore
		return nil
	}
}

// OptIsCA sets the is certificate authority flag.
func OptIsCA(isCA bool) CertOption {
	return func(csr *CertOptions) error {
		csr.IsCA = isCA
		return nil
	}
}

// OptKeyUsage sets the key usage flags.
func OptKeyUsage(keyUsage x509.KeyUsage) CertOption {
	return func(csr *CertOptions) error {
		csr.KeyUsage = keyUsage
		return nil
	}
}

// OptDNSNames sets valid dns names for the cert.
func OptDNSNames(dnsNames ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.DNSNames = dnsNames
		return nil
	}
}

// OptAddDNSNames adds valid dns names for the cert.
func OptAddDNSNames(dnsNames ...string) CertOption {
	return func(csr *CertOptions) error {
		csr.DNSNames = append(csr.DNSNames, dnsNames...)
		return nil
	}
}

// OptSerialNumber sets the serial number for the certificate.
// If this option isn't provided, a random one is generated.
func OptSerialNumber(serialNumber *big.Int) CertOption {
	return func(cco *CertOptions) error {
		cco.SerialNumber = serialNumber
		return nil
	}
}

// OptPrivateKey sets the private key to use when generating the certificate.
// If this option isn't provided, a new one is generated.
func OptPrivateKey(privateKey *rsa.PrivateKey) CertOption {
	return func(cco *CertOptions) error {
		cco.PrivateKey = privateKey
		return nil
	}
}

// OptPrivateKeyFromPath reads a private key from a given path and parses it as PKCS1PrivateKey.
func OptPrivateKeyFromPath(path string) CertOption {
	return func(cco *CertOptions) error {
		contents, err := os.ReadFile(path)
		if err != nil {
			return ex.New(err)
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(contents)
		if err != nil {
			return ex.New(err)
		}
		cco.PrivateKey = privateKey
		return nil
	}
}
