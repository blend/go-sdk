/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
	_ graceful.Graceful = (*GracefulHTTPServer)(nil)
)

func TestGracefulServer(t *testing.T) {
	assert := assert.New(t)

	listener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	typedListener, ok := listener.(*net.TCPListener)
	assert.True(ok)
	assert.NotNil(typedListener)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK!\n")
		}),
	}
	gs := NewGracefulHTTPServer(server, OptGracefulHTTPServerListener(typedListener))
	stopSignal := make(chan os.Signal)
	didShutdown := make(chan struct{})

	go func() {
		defer func() { close(didShutdown) }()
		_ = graceful.ShutdownBySignal([]graceful.Graceful{gs}, graceful.OptShutdownSignal(stopSignal))
	}()
	<-gs.NotifyStarted()

	res, err := http.Get("http://" + typedListener.Addr().String())
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	stopSignal <- os.Interrupt
	<-didShutdown
}
