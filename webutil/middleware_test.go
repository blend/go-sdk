/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMiddelware(t *testing.T) {
	assert := assert.New(t)

	innerDone := make(chan struct{})
	inner := func(rw http.ResponseWriter, req *http.Request) {
		close(innerDone)
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, "OK!\n")
	}

	oneDone := make(chan struct{})
	one := func(action http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			close(oneDone)
			action(rw, req)
			fmt.Fprint(rw, "One\n")
		}
	}

	twoDone := make(chan struct{})
	two := func(action http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			close(twoDone)
			action(rw, req)
			fmt.Fprint(rw, "Two\n")
		}
	}

	server := httptest.NewServer(http.HandlerFunc(NestMiddleware(inner, two, one)))
	defer server.Close()

	res, err := http.Get(server.URL)
	assert.Nil(err)

	<-oneDone
	<-twoDone
	<-innerDone

	contents, err := io.ReadAll(res.Body)
	assert.Nil(err)
	assert.Equal("OK!\nTwo\nOne\n", string(contents))
}
