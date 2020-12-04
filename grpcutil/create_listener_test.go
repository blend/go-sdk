package grpcutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestListener(t *testing.T) {
	assert := assert.New(t)

	tcpln, err := CreateListener("127.0.0.1:")
	assert.Nil(err)
	defer tcpln.Close()
	assert.Equal("tcp", tcpln.Addr().Network())
	assert.Contains(tcpln.Addr().String(), "127.0.0.1:")

	socketPath := fmt.Sprintf("/tmp/%s.sock", uuid.V4().String())
	socketAddress := fmt.Sprintf("unix://" + socketPath)
	unixln, err := CreateListener(socketAddress)
	assert.Nil(err)
	defer unixln.Close()
	assert.Equal("unix", unixln.Addr().Network())
	assert.Equal(socketPath, unixln.Addr().String())
}
