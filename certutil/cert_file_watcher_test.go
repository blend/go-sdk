/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"crypto/x509"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCertFileWatcher(t *testing.T) {
	its := assert.New(t)

	tempDir, err := os.MkdirTemp("", "")
	its.Nil(err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	tempCertPath := filepath.Join(tempDir, "tls.crt")
	tempKeyPath := filepath.Join(tempDir, "tls.key")

	err = copyFile("testdata/server.cert.pem", tempCertPath)
	its.Nil(err)
	err = copyFile("testdata/server.key.pem", tempKeyPath)
	its.Nil(err)

	w, err := NewCertFileWatcher(
		KeyPair{CertPath: tempCertPath, KeyPath: tempKeyPath},
	)
	its.Nil(err)

	its.Equal(tempCertPath, w.CertPath())
	its.Equal(tempKeyPath, w.KeyPath())

	cert := w.Certificate()
	its.NotNil(cert)

	err = copyFile("testdata/alt-server.cert.pem", tempCertPath)
	its.Nil(err)
	err = copyFile("testdata/alt-server.key.pem", tempKeyPath)
	its.Nil(err)

	err = w.Reload()
	its.Nil(err)

	newCert := w.Certificate()
	its.NotNil(newCert)

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	its.Nil(err)
	newCert.Leaf, err = x509.ParseCertificate(newCert.Certificate[0])
	its.Nil(err)

	its.NotEqual(cert.Leaf.SerialNumber.String(), newCert.Leaf.SerialNumber.String())
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
