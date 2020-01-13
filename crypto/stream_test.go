package crypto

import (
	"bytes"
	"io/ioutil"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestStreamEncrypterDecrypter(t *testing.T) {
	assert := assert.New(t)
	key, err := CreateKey(32)
	assert.Nil(err)
	plaintext := "Eleven is the best person in all of Hawkins Indiana. Some more text"
	pt := []byte(plaintext)

	src := bytes.NewReader(pt)

	se, err := NewStreamEncrypter(key, src)
	assert.Nil(err)
	assert.NotNil(se)

	encrypted, err := ioutil.ReadAll(se)
	assert.Nil(err)
	assert.NotNil(encrypted)

	sd, err := NewStreamDecrypter(key, se.Meta(), bytes.NewReader(encrypted))
	assert.Nil(err)
	assert.NotNil(sd)

	decrypted, err := ioutil.ReadAll(sd)
	assert.Nil(err)
	assert.Equal(plaintext, string(decrypted))

	assert.Nil(sd.Authenticate())
}

func TestCheckedWrite(t *testing.T) {
	assert := assert.New(t)
	writer := bytes.NewBuffer(nil)
	data := []byte{1, 1, 1}
	v, err := checkedWrite(writer, data)
	assert.Nil(err)
	assert.Equal(len(data), v)
}
