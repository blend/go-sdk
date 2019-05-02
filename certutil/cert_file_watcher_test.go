package certutil

import (
	"crypto/tls"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCertFileWatcher(t *testing.T) {
	assert := assert.New(t)

	tempKey, err := ioutil.TempFile("", "")
	assert.Nil(err)
	defer func() {
		os.Remove(tempKey.Name())
	}()

	tempCert, err := ioutil.TempFile("", "")
	assert.Nil(err)
	defer func() {
		os.Remove(tempCert.Name())
	}()

	_, err = tempKey.Write(keyLiteral)
	assert.Nil(err)

	_, err = tempCert.Write(certLiteral)
	assert.Nil(err)

	assert.Nil(tempKey.Close())
	assert.Nil(tempCert.Close())

	w, err := NewCertFileWatcher(tempCert.Name(), tempKey.Name())
	assert.Nil(err)
	assert.NotNil(w.Certificate)

	assert.Nil(w.Reload())
	assert.NotNil(w.Certificate)

	reloaded := make(chan struct{})
	w.OnReload = func(_ *tls.Certificate, _ error) {
		close(reloaded)
	}

	w.PollInterval = 5 * time.Millisecond

	go w.Start()
	<-w.NotifyStarted()
	defer w.Stop()

	kw, err := os.OpenFile(tempKey.Name(), os.O_RDWR, 0644)
	assert.Nil(err)

	cw, err := os.OpenFile(tempCert.Name(), os.O_RDWR, 0644)
	assert.Nil(err)

	_, err = kw.WriteAt(alternateKeyLiteral, 0)
	assert.Nil(err)
	assert.Nil(kw.Close())

	_, err = cw.WriteAt(alternateCertLiteral, 0)
	assert.Nil(err)
	assert.Nil(cw.Close())

	<-reloaded
	assert.NotNil(w.Certificate)
}
