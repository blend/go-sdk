package proxyprotocol

import (
	"crypto/tls"
	"net"
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
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

func TestCreateTLSListener(t *testing.T) {
	assert := assert.New(t)

	tlsConfig := &tls.Config{}
	listener, err := CreateListener("127.0.0.1:",
		OptKeepAlive(true),
		OptUseProxyProtocol(true),
		OptKeepAlivePeriod(30*time.Second),
		OptTLSConfig(tlsConfig),
	)
	defer listener.Close()

	assert.Nil(err)
	assert.NotNil(listener)

	go func() {
		_, err := net.Dial("tcp", listener.Addr().String())
		assert.Nil(err)
	}()

	conn, err := listener.Accept()
	assert.Nil(err)

	typed, ok := conn.(*tls.Conn)
	assert.True(ok)
	assert.NotNil(typed)
}
