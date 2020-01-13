package certutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestExtendSystemCertPool(t *testing.T) {
	assert := assert.New(t)

	pool, err := ExtendSystemCertPool(KeyPair{Cert: string(caCertLiteral)})
	assert.Nil(err)
	assert.NotNil(pool)
}
