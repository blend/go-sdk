package request

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestNewCertInfo(t *testing.T) {
	assert := assert.New(t)

	notBefore, notAfter := time.Now(), time.Now()
	issuer := "go-sdk tests"
	subject := "test.blend.com"
	dnsNames := []string{subject}
	response := &http.Response{
		TLS: &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{
				&x509.Certificate{
					Issuer: pkix.Name{
						CommonName: issuer,
					},
					Subject: pkix.Name{
						CommonName: subject,
					},
					DNSNames:  dnsNames,
					NotBefore: notBefore,
					NotAfter:  notAfter,
				},
			},
		},
	}

	// Test: NewCertInfo should include dns names, validity periods, issuer and subject common names
	certInfo := NewCertInfo(response)
	assert.Equal(dnsNames, certInfo.DNSNames)
	assert.Equal(notBefore, certInfo.NotBefore)
	assert.Equal(notAfter, certInfo.NotAfter)
	assert.Equal(issuer, certInfo.IssuerCommonName)
	assert.Equal(subject, certInfo.SubjectCommonName)
}
