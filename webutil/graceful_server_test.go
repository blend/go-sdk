package webutil

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

var (
	_ graceful.Graceful = (*GracefulServer)(nil)
)

func TestGracefulServer(t *testing.T) {
	assert := assert.New(t)

	listener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	typedListener, ok := listener.(*net.TCPListener)
	assert.True(ok)

	server := &http.Server{}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	})
	gs := NewGracefulServer(server).WithListener(typedListener)
	stopSignal := make(chan os.Signal)
	didShutdown := make(chan struct{})
	go func() {
		defer func() { close(didShutdown) }()
		graceful.ShutdownBySignal(gs, stopSignal)
	}()
	res, err := http.Get("http://" + typedListener.Addr().String())
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	stopSignal <- os.Interrupt
	<-didShutdown
}
