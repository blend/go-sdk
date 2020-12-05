package webutil

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestParseCertInfo(t *testing.T) {
	assert := assert.New(t)

	// handle the empty cases
	assert.Nil(ParseCertInfo(nil))
	assert.Nil(ParseCertInfo(&http.Response{}))
	assert.Nil(ParseCertInfo(&http.Response{
		TLS: &tls.ConnectionState{},
	}))

	na := time.Now().UTC().AddDate(0, 1, 0)
	nb := time.Now().UTC().AddDate(0, -1, 0)
	valid := &http.Response{
		TLS: &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{
				&x509.Certificate{
					Issuer: pkix.Name{
						CommonName: "bailey dog",
						Names: []pkix.AttributeTypeAndValue{
							pkix.AttributeTypeAndValue{Value: "x"},
							pkix.AttributeTypeAndValue{Value: "yz"},
							pkix.AttributeTypeAndValue{Value: "abc"},
						},
					},
					DNSNames:  []string{"foo.local"},
					NotAfter:  na,
					NotBefore: nb,
				},
			},
		},
	}

	info := ParseCertInfo(valid)
	expected := &CertInfo{
		IssuerCommonName: "bailey dog",
		IssuerNames:      []string{"x", "yz", "abc"},
		DNSNames:         []string{"foo.local"},
		NotAfter:         na,
		NotBefore:        nb,
	}
	assert.Equal(expected, info)

	assert.False(info.IsExpired())
	assert.False(info.WillBeExpired(time.Now().UTC()))
}

func TestCertInfoWillBeExpired(t *testing.T) {
	assert := assert.New(t)

	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	t2 := time.Date(2009, time.November, 10, 23, 0, 1, 0, time.UTC)

	info := &CertInfo{}
	// Both `NotAfter` and `NotBefore` un set.
	assert.False(info.WillBeExpired(t1))
	// Only `NotAfter` set, will not be expired.
	info = &CertInfo{NotAfter: t2}
	assert.False(info.WillBeExpired(t1))
	// Only `NotAfter` set, will be expired.
	info = &CertInfo{NotAfter: t1}
	assert.True(info.WillBeExpired(t2))
	// Only `NotBefore` set, will not be expired.
	info = &CertInfo{NotBefore: t1}
	assert.False(info.WillBeExpired(t2))
	// Only `NotBefore` set, will be expired.
	info = &CertInfo{NotBefore: t2}
	assert.True(info.WillBeExpired(t1))
}
