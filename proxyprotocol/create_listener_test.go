package proxyprotocol

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/webutil"

	"github.com/blend/go-sdk/assert"
)

func TestCreateListener(t *testing.T) {
	assert := assert.New(t)

	listener, err := CreateListener("127.0.0.1:",
		OptKeepAlive(true),
		OptUseProxyProtocol(true),
		OptKeepAlivePeriod(30*time.Second),
	)
	defer listener.Close()

	assert.Nil(err)
	assert.NotNil(listener)

	typed, ok := listener.(*Listener)
	assert.True(ok)
	assert.NotNil(typed)

	assert.NotNil(typed.Listener)

	tcpListener, ok := typed.Listener.(webutil.TCPKeepAliveListener)
	assert.True(ok)
	assert.NotNil(tcpListener)
}
