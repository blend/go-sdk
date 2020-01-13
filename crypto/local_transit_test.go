package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestLocalTransitEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)

	plaintext := "mary jane hawkins"

	m := NewLocalTransit(OptLocalTransitContextProvider(func() string {
		return time.Date(2019, 04, 15, 01, 02, 03, 04, time.UTC).Format("20060102")
	}), OptLocalTransitKey(MustCreateKey(32)))
	prefix := m.ContextProvider()

	ciphertext := new(bytes.Buffer)
	assert.Nil(m.Encrypt(ciphertext, bytes.NewReader([]byte(plaintext))))

	cipherBytes := ciphertext.Bytes()

	assert.True(len(cipherBytes) > KeyVersionSize+IVSize+HashSize)
	assert.True(strings.HasPrefix(string(cipherBytes), prefix+":"), "we should prefix ciphertext with the current date")

	output := new(bytes.Buffer)
	assert.Nil(m.Decrypt(output, bytes.NewReader(ciphertext.Bytes())))

	assert.Equal(plaintext, output.String())
}

func TestEncryptDecryptLarge(t *testing.T) {
	assert := assert.New(t)

	m := NewLocalTransit(OptLocalTransitKey(MustCreateKey(32)))
	m.ContextProvider = func() string {
		return time.Date(2019, 04, 15, 01, 02, 03, 04, time.UTC).Format("20060102")
	}
	plaintext := make([]byte, 64*1024) // 64kb of data
	_, err := rand.Read(plaintext)
	assert.Nil(err)

	ciphertext := new(bytes.Buffer)
	assert.Nil(m.Encrypt(ciphertext, bytes.NewReader(plaintext)))

	output := new(bytes.Buffer)
	assert.Nil(m.Decrypt(output, ciphertext))
	assert.True(hmac.Equal(output.Bytes(), plaintext))
}
