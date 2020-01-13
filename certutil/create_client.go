package certutil

import (
	"crypto/rand"
	"crypto/x509"
	"github.com/blend/go-sdk/ex"
)

/*
CreateClient creates a client cert bundle associated with a given common name.

The CA must be passed in as a CertBundle.

Example:

	ca, err := certutil.NewCertBundle(certutil.KeyPairFromPaths("ca.crt", "ca.key"))
	if err != nil {
		return err
	}
	client, err := CreateClient("foo.bar.com", ca)
*/
func CreateClient(commonName string, ca *CertBundle, options ...CertOption) (*CertBundle, error) {
	if ca == nil {
		return nil, ex.New("must provide a ca cert bundle")
	}

	createOptions := DefaultOptionsClient
	createOptions.Subject.CommonName = commonName
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
