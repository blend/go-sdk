/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSClientCertFilePair(t *testing.T) {
	assert := assert.New(t)

	certFile, err := writeTempFile(clientCert)
	assert.Nil(err)
	defer os.Remove(certFile)
	keyFile, err := writeTempFile(clientKey)
	assert.Nil(err)
	defer os.Remove(keyFile)

	r := New("https://foo.com", OptTLSClientCertFilePair(certFile, keyFile))
	assert.NotNil(r.Client)
	assert.NotNil(r.Client.Transport)
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig)
	assert.NotEmpty(r.Client.Transport.(*http.Transport).TLSClientConfig.Certificates)
}

func TestOptTLSClientCertFilePairErrors(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.com", OptTLSClientCertFilePair("", ""))
	assert.NotNil(r.Err)
}

func writeTempFile(contents []byte) (string, error) {
	tf, err := os.CreateTemp("", "r2_opt_tls_client_cert_file_pair")
	if err != nil {
		return "", err
	}
	defer tf.Close()
	_, err = io.Copy(tf, bytes.NewReader(contents))
	if err != nil {
		return "", err
	}
	return tf.Name(), nil
}
