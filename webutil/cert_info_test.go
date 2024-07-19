/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

	valid := &http.Response{
		TLS: &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{
				{
					Issuer: pkix.Name{
						CommonName: "example-string dog",
					},
					DNSNames:  []string{"foo.local"},
					NotAfter:  time.Now().UTC().AddDate(0, 1, 0),
					NotBefore: time.Now().UTC().AddDate(0, -1, 0),
				},
			},
		},
	}

	info := ParseCertInfo(valid)
	assert.NotNil(info)
	assert.Equal("example-string dog", info.IssuerCommonName)
	assert.Equal([]string{"foo.local"}, info.DNSNames)
	assert.False(info.NotAfter.IsZero())
	assert.False(info.NotBefore.IsZero())
	assert.True(info.NotAfter.After(time.Now().UTC()))
	assert.True(info.NotBefore.Before(time.Now().UTC()))

	assert.False(info.IsExpired())
	assert.False(info.WillBeExpired(time.Now().UTC()))
}
