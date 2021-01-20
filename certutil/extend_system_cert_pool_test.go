package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestExtendSystemCertPool(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	pool, err := ExtendSystemCertPool(KeyPair{Cert: string(caCertLiteral)})
	assert.Nil(err)
	assert.NotNil(pool)
}
