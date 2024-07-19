/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/blend/go-sdk/certutil"
)

const testCertChain = `-----BEGIN CERTIFICATE-----
MIIDBDCCAeygAwIBAgIRAJmlkpC8Q9rL5W+8tNKDaz0wDQYJKoZIhvcNAQELBQAw
ETEPMA0GA1UEChMGZ28tc2RrMB4XDTE5MDMyOTE3MTMxM1oXDTI0MDMyOTE3MTMx
M1owKzEVMBMGA1UEChMMZ28tc2RrIHVzZXJzMRIwEAYDVQQDEwlkZXYubG9jYWww
ggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDQNLHSwaKJCBUHSP5MwQ/R
f+b2Lzh+DLN+CAPkoIWxZ99cjc+hUJ2WvdPmMzq4jfk+KhWzLWps45iR+jSFV5GQ
gxI8ZmyLH9QwA5byr22Sy3P7XNbebUGz0AIU3ZxgQLtWRhC9bZ0vnLrmK/cwaEAz
Q6Jh4YKOSafgiL8wb7QfieKhdkgwsNctxSmDe//V9pyiFfDtcECQA4bzb+XHQn6+
+w8FlERFvVmLRUiR29jvbmbPCJF/VH/244KUztzap9BVIbMPgvjYya615Mou6iGE
eomAue/XN01m37ky1k+C6PLzFqoOIWI2bw+pAq+GH7ijt2n2nevuT8yvWnsnY6IN
AgMBAAGjPTA7MA4GA1UdDwEB/wQEAwIHgDATBgNVHSUEDDAKBggrBgEFBQcDATAU
BgNVHREEDTALgglkZXYubG9jYWwwDQYJKoZIhvcNAQELBQADggEBACDdtI3LRnRR
57YGXJ4Pyw10V/W84qNk2uWPEsMQESO8ufc8uS1vqLER5EIuKfbgAe/+U3aX0wyZ
sCr39kqFOITq7yfBwm0Ige81n5KVy+hOlGDdkG37ol/bGB1+rf7MClokmAVFhtI2
kfeWSfLI+SODvtvM7nnpyCwtvJjWNILcoEl6BKByHQ3vbDEKs2w4XGsftfyy5vGW
YABlakh47OldjuGzSAPJUECsOOK9RlrZLynNqiGRfj/LaS2C0zZrDM89kuPngM9U
X+maxcJQFXbK7NT+SoiVeT54vTbwSm3friwRkh2JTRnAZKhDSCHJEsgVXVh2XMjH
LXfRlJqO8qk=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIC0DCCAbigAwIBAgIRAMA5rXAuCDKJ4DUJi5TffXUwDQYJKoZIhvcNAQELBQAw
ETEPMA0GA1UEChMGZ28tc2RrMB4XDTE5MDMyOTE3MTMxM1oXDTE5MDQyODE3MTMx
M1owETEPMA0GA1UEChMGZ28tc2RrMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAp4XHxZEPKoRbHSOLSPIMN1XoMlvEFYbfkwkZEKCpVw2pBZIB/9KD8JLF
aMbsYxouPjPLq10t460qQWNVBKQACeq9hUOHN2ZtV0oH6a//fdxqDkinrWIhhVeg
MSxWCIYqUlyjiSdFXTCIjuI2FNKN9S2fcvCeLqGjhU3IFypG9dYKZkQ2Oqnejbuo
TLL6wz85T9MN5cr2gbAmWGAu5SFn8gEeJImgxbRMPebFejk6X1hwgzpVicgXE3Ib
ZxxCJKfLfgwfEhvqd6JC+A4UgfcnqMNTic/ztkEJDdhq6OPkKaoanlZUlamBOuwt
T8rqsa8ekHnJjRNDfwKPvM0o7i6XxwIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAqQw
DwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEANE8vXN1s+ohBi3+S
xnnSr3CCamS6h9MSRmEb1Q1jCZ+HQLUeQgIdLCbwxPB10i4fJkV5vTo9r/OER/G5
XhAKT1IZexLhcbyr9tf3s7Xd43QU1MU4Q3gRHra5w+RNgi3Dnq1zp6OHwdPDUA9R
INM2BEhXOwDefHYmEHATCCeyr3GGm12Ot1EdFfy6YCU5Q6xiHB4/liTmHcbruYAP
XMNCo8B1Pd/psRLF8hcqp0jeHvy54GfJwT1WvHzdpQ+yB5TbsVyxTH3jI/E4s6ZL
LMb3TkFNROSPTD9CVPU1kVVCI7A2X/AKtEc1sbmv0G2FLH2UqzsLqJKj8ZSmDfGP
H/db8w==
-----END CERTIFICATE-----
`

const days = 24 * time.Hour

func main() {
	warningThreshold := 60 * days

	parsedCerts, err := certutil.ParseCertPEM([]byte(testCertChain))
	if err != nil {
		log.Fatal(err)
	}
	if len(parsedCerts) == 0 {
		log.Fatal(errors.New("invalid pem; empty certs"))
	}
	for _, cert := range parsedCerts {
		if time.Now().UTC().After(cert.NotAfter) {
			fmt.Printf("fatal; certificate %s expired %s", cert.Subject.CommonName, cert.NotAfter.Format(time.RFC3339))
		} else if delta := cert.NotAfter.Sub(time.Now().UTC()); delta < warningThreshold {
			fmt.Printf("warning; certificate %s will expire in %v", cert.Subject.CommonName, formatDays(delta))
		} else {
			fmt.Printf("ok; certificate %s is still valid", cert.Subject.CommonName)
		}
		fmt.Println()
	}
}

func formatDays(d time.Duration) string {
	if d > days {
		return fmt.Sprintf("%d days", d/days)
	}
	return d.String()
}
