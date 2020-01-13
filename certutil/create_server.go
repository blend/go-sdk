package certutil

import (
	"crypto/rand"
	"crypto/x509"
	"github.com/blend/go-sdk/ex"
)

// CreateServer creates a ca cert bundle.
func CreateServer(commonName string, ca *CertBundle, options ...CertOption) (*CertBundle, error) {
	if ca == nil || ca.PrivateKey == nil || len(ca.Certificates) == 0 {
		return nil, ex.New("provided certificate authority bundle is invalid")
	}

	createOptions := DefaultOptionsServer
	// set the default common name
	createOptions.Subject.CommonName = commonName
	// it is important to reflect the common name here as well
	createOptions.DNSNames = []string{commonName}

	if err := ResolveCertOptions(&createOptions, options...); err != nil {
		return nil, nil
	}

	var output CertBundle
	output.PrivateKey = createOptions.PrivateKey
	output.PublicKey = &createOptions.PrivateKey.PublicKey
	der, err := x509.CreateCertificate(rand.Reader, &createOptions.Certificate, &ca.Certificates[0], output.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, ex.New(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, ex.New(err)
	}
	output.CertificateDERs = append([][]byte{der}, ca.CertificateDERs...)
	output.Certificates = append([]x509.Certificate{*cert}, ca.Certificates...)
	return &output, nil
}
