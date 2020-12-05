package webutil

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTCPKeepAliveListenerAccept(t *testing.T) {
	assert := assert.New(t)

	var wg sync.WaitGroup
	wg.Add(1)
	port := 10000 + rand.Intn(0xffff-10000)
	addr := fmt.Sprintf("http://127.0.0.1:%d", port)
	var connectResp *http.Response
	var connectErr error
	go func() {
		defer wg.Done()
		connectResp, connectErr = http.Get(addr)
	}()

	ln, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	assert.Nil(err)

	tkal := TCPKeepAliveListener{
		TCPListener:     ln,
		KeepAlive:       true,
		KeepAlivePeriod: 3 * time.Millisecond,
	}
	c, err := tkal.Accept()
	assert.Nil(err)
	assert.NotNil(c)
	// NOTE: Unfortunately there is no simple way to verify that
	//       `SetKeepAlive` / `SetKeepAlivePeriod` were called.

	// Clean up the request.
	assert.Nil(c.Close())
	wg.Wait()
	assert.Nil(connectResp)
	assert.NotNil(connectErr)
	unwrapped, ok := connectErr.(*url.Error)
	assert.True(ok)
	// NOTE: `unwrapped.Err` will vary across runs, so we do not check.
	assert.Equal("Get", unwrapped.Op)
	assert.Equal(addr, unwrapped.URL)

	// Accept() fails.
	port = 10000 + rand.Intn(0xffff-10000)
	ln, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	assert.Nil(err)
	assert.NotNil(ln)
	assert.Nil(ln.Close())

	tkal = TCPKeepAliveListener{TCPListener: ln}
	// NOTE: A closed listener cannot accept a connection.
	c, err = tkal.Accept()
	assert.Nil(c)
	assert.NotNil(err)
	expected := fmt.Sprintf("accept tcp 127.0.0.1:%d: use of closed network connection", port)
	assert.Equal(expected, err.Error())
}
