package crypto

import (
	"crypto/hmac"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestCreateKey(t *testing.T) {
	assert := assert.New(t)

	key, err := CreateKey(32)
	assert.Nil(err)
	assert.Len(key, 32)

	key2, err := CreateKey(32)
	assert.Nil(err)
	assert.Len(key2, 32)

	assert.False(hmac.Equal(key, key2))
}
