package webutil

import (
	"fmt"
	"net/http"
	"time"
)

// ParseCertInfo returns a new cert info from a response from a check.
func ParseCertInfo(res *http.Response) *CertInfo {
	if res == nil || res.TLS == nil || len(res.TLS.PeerCertificates) == 0 {
		return nil
	}

	var earliestNotAfter time.Time
	var latestNotBefore time.Time
	for _, cert := range res.TLS.PeerCertificates {
		if earliestNotAfter.IsZero() || earliestNotAfter.After(cert.NotAfter) {
			earliestNotAfter = cert.NotAfter
		}
		if latestNotBefore.IsZero() || latestNotBefore.Before(cert.NotBefore) {
			latestNotBefore = cert.NotBefore
		}
	}

	firstCert := res.TLS.PeerCertificates[0]
	var issuerNames []string
	for _, name := range firstCert.Issuer.Names {
		issuerNames = append(issuerNames, fmt.Sprint(name.Value))
	}

	return &CertInfo{
		IssuerCommonName: firstCert.Issuer.CommonName,
		IssuerNames:      issuerNames,
		DNSNames:         firstCert.DNSNames,
		NotAfter:         earliestNotAfter,
		NotBefore:        latestNotBefore,
	}
}

// CertInfo is the information for a certificate.
type CertInfo struct {
	IssuerCommonName string    `json:"issuerCommonName" yaml:"issuerCommonName"`
	IssuerNames      []string  `json:"issuerNames" yaml:"issuerNames"`
	DNSNames         []string  `json:"dnsNames" yaml:"dnsNames"`
	NotAfter         time.Time `json:"notAfter" yaml:"notAfter"`
	NotBefore        time.Time `json:"notBefore" yaml:"notBefore"`
}

// IsExpired returns if the certificate is strictly expired
// and would not be accepted by browsers.
func (ci CertInfo) IsExpired() bool {
	return ci.WillBeExpired(time.Now().UTC())
}

// WillBeExpired returns if the certificate is strictly expired
// and would not be accepted by browsers at a given time.
func (ci CertInfo) WillBeExpired(at time.Time) bool {
	if !ci.NotAfter.IsZero() {
		if at.UTC().After(ci.NotAfter) {
			return true
		}
	}
	if !ci.NotBefore.IsZero() {
		if at.UTC().Before(ci.NotBefore) {
			return true
		}
	}
	return false
}
