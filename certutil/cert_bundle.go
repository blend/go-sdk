package certutil

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"github.com/blend/go-sdk/ex"
)

// NewCertBundle returns a new cert bundle from a given key pair, which can denote the raw PEM encoded
// contents of the public and private key portions of the cert, or paths to files.
// The CertBundle itself is the parsed public key, private key, and individual certificates for the pair.
func NewCertBundle(keyPair KeyPair) (*CertBundle, error) {
	certPEM, err := keyPair.CertBytes()
	if err != nil {
		return nil, ex.New(err)
	}
	if len(certPEM) == 0 {
		return nil, ex.New("empty cert contents")
	}

	keyPEM, err := keyPair.KeyBytes()
	if err != nil {
		return nil, ex.New(err)
	}
	if len(keyPEM) == 0 {
		return nil, ex.New("empty key contents")
	}

	certData, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, ex.New(err)
	}
	if len(certData.Certificate) == 0 {
		return nil, ex.New("no certificates")
	}

	var certs []x509.Certificate
	var ders [][]byte
	for _, certDataPortion := range certData.Certificate {
		cert, err := x509.ParseCertificate(certDataPortion)
		if err != nil {
			return nil, ex.New(err)
		}

		certs = append(certs, *cert)
		ders = append(ders, cert.Raw)
	}

	var privateKey *rsa.PrivateKey
	if typed, ok := certData.PrivateKey.(*rsa.PrivateKey); ok {
		privateKey = typed
	} else {
		return nil, ex.New("invalid private key type", ex.OptMessagef("%T", certData.PrivateKey))
	}

	return &CertBundle{
		PrivateKey:      privateKey,
		PublicKey:       &privateKey.PublicKey,
		Certificates:    certs,
		CertificateDERs: ders,
	}, nil
}

// CertBundle is the packet of information for a certificate.
type CertBundle struct {
	PrivateKey      *rsa.PrivateKey
	PublicKey       *rsa.PublicKey
	Certificates    []x509.Certificate
	CertificateDERs [][]byte
}

// MustGenerateKeyPair returns a serialized version of the bundle as a key pair
// and panics if there is an error.
func (cb *CertBundle) MustGenerateKeyPair() KeyPair {
	pair, err := cb.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	return pair
}

// GenerateKeyPair returns a serialized key pair for the cert bundle.
func (cb *CertBundle) GenerateKeyPair() (output KeyPair, err error) {
	private := bytes.NewBuffer(nil)
	if err = cb.WriteKeyPem(private); err != nil {
		return
	}
	public := bytes.NewBuffer(nil)
	if err = cb.WriteCertPem(public); err != nil {
		return
	}
	output = KeyPair{
		Cert: public.String(),
		Key:  private.String(),
	}
	return
}

// WithParent adds a parent certificate to the certificate chain.
// It is used typically to add the certificate authority.
func (cb *CertBundle) WithParent(parent *CertBundle) {
	cb.Certificates = append(cb.Certificates, parent.Certificates...)
	cb.CertificateDERs = append(cb.CertificateDERs, parent.CertificateDERs...)
}

// WriteCertPem writes the public key portion of the cert to a given writer.
func (cb CertBundle) WriteCertPem(w io.Writer) error {
	for _, der := range cb.CertificateDERs {
		if err := pem.Encode(w, &pem.Block{Type: BlockTypeCertificate, Bytes: der}); err != nil {
			return ex.New(err)
		}
	}
	return nil
}

// CertPEM returns the cert portion of the certificate DERs as a byte array.
func (cb CertBundle) CertPEM() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := cb.WriteCertPem(buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// WriteKeyPem writes the certificate key as a pem.
func (cb CertBundle) WriteKeyPem(w io.Writer) error {
	return pem.Encode(w, &pem.Block{Type: BlockTypeRSAPrivateKey, Bytes: x509.MarshalPKCS1PrivateKey(cb.PrivateKey)})
}

// KeyPEM returns the cert portion of the certificate DERs as a byte array.
func (cb CertBundle) KeyPEM() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := cb.WriteKeyPem(buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// CommonNames returns the cert bundle common name(s).
func (cb CertBundle) CommonNames() ([]string, error) {
	if len(cb.Certificates) == 0 {
		return nil, ex.New("no certificates returned")
	}
	var output []string
	for _, cert := range cb.Certificates {
		if cert.Subject.CommonName != "" {
			output = append(output, cert.Subject.CommonName)
		}
	}
	return output, nil
}

// CertPool returns the bundle as a cert pool.
func (cb CertBundle) CertPool() (*x509.CertPool, error) {
	systemPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, ex.New(err)
	}
	for index := range cb.Certificates {
		systemPool.AddCert(&cb.Certificates[index])
	}
	return systemPool, nil
}
