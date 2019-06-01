package grpcutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestListener(t *testing.T) {
	assert := assert.New(t)

	tcpln, err := Listener("127.0.0.1:8080")
	assert.Nil(err)
	defer tcpln.Close()
	assert.Equal("tcp", tcpln.Addr().Network())
	assert.Equal("127.0.0.1:8080", tcpln.Addr().String())

	unixln, err := Listener("unix:///tmp/test.sock")
	assert.Nil(err)
	defer unixln.Close()
	assert.Equal("unix", unixln.Addr().Network())
	assert.Equal("/tmp/test.sock", unixln.Addr().String())
}
