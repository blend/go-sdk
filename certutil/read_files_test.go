package certutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestReadFiles(t *testing.T) {
	assert := assert.New(t)

	files, err := ReadFiles("testdata/client.cert.pem", "testdata/client.key.pem")
	assert.Nil(err)
	assert.Len(files, 2)
}
