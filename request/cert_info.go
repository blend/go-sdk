package request

import (
	"fmt"
	"net/http"
	"time"
)

// NewCertInfo returns a new cert info from a response.
func NewCertInfo(res *http.Response) *CertInfo {
	if res.TLS != nil && len(res.TLS.PeerCertificates) > 0 {
		var earliestExpiration time.Time
		var latestNotBefore time.Time
		for _, cert := range res.TLS.PeerCertificates {
			if earliestExpiration.IsZero() || earliestExpiration.After(cert.NotAfter) {
				earliestExpiration = cert.NotAfter
			}
			if latestNotBefore.IsZero() || latestNotBefore.Before(cert.NotBefore) {
				latestNotBefore = cert.NotBefore
			}
		}

		firstCert := res.TLS.PeerCertificates[0]

		var issuerCommonName string
		if len(firstCert.Issuer.CommonName) > 0 {
			issuerCommonName = firstCert.Issuer.CommonName
		} else {
			for _, name := range firstCert.Issuer.Names {
				if name.Type.String() == "2.5.4.3" {
					issuerCommonName = fmt.Sprintf("%v", name.Value)
				}
			}
		}

		return &CertInfo{
			DNSNames:          firstCert.DNSNames,
			NotAfter:          earliestExpiration,
			NotBefore:         latestNotBefore,
			IssuerCommonName:  issuerCommonName,
			SubjectCommonName: firstCert.Subject.CommonName,
		}
	}

	return nil
}

// CertInfo is the information for a certificate.
type CertInfo struct {
	IssuerCommonName  string    `json:"issuerCommonName" yaml:"issuerCommonName"`
	SubjectCommonName string    `json:"subjectCommonName" yaml:"subjectCommonName"`
	DNSNames          []string  `json:"dnsNames" yaml:"dnsNames"`
	NotAfter          time.Time `json:"notAfter" yaml:"notAfter"`
	NotBefore         time.Time `json:"notBefore" yaml:"notBefore"`
}
