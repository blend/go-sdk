package certutil

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCertWatcher(t *testing.T) {
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

	ca, err := CreateCertificateAuthority()
	assert.Nil(err)

	server, err := CreateServer("local.test", ca)
	assert.Nil(err)

	pair, err := server.GenerateKeyPair()
	_, err = tempKey.Write([]byte(pair.Key))
	assert.Nil(err)

	_, err = tempCert.Write([]byte(pair.Cert))
	assert.Nil(err)

	assert.Nil(tempKey.Close())
	assert.Nil(tempCert.Close())

	w, err := NewCertWatcher(tempCert.Name(), tempKey.Name())
	assert.Nil(err)
	assert.NotNil(w.Certificate)

	assert.Nil(w.Reload())
	assert.NotNil(w.Certificate)

	reloaded := make(chan struct{})
	w.OnReload = func() {
		close(reloaded)
	}

	w.PollInterval = time.Microsecond
	go w.Start()
	defer w.Stop()

	// recreate the server ...
	server, err = CreateServer("local.test", ca)
	assert.Nil(err)

	pair, err = server.GenerateKeyPair()

	kw, err := os.OpenFile(tempKey.Name(), os.O_RDWR, 0644)
	assert.Nil(err)
	_, err = kw.WriteAt([]byte(pair.Key), 0)
	assert.Nil(err)
	assert.Nil(kw.Close())

	cw, err := os.OpenFile(tempCert.Name(), os.O_RDWR, 0644)
	assert.Nil(err)
	_, err = cw.WriteAt([]byte(pair.Cert), 0)
	assert.Nil(err)
	assert.Nil(cw.Close())

	<-reloaded
	assert.NotNil(w.Certificate)
}
