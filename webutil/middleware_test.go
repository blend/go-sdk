package webutil

import (
	"io"
	"io/ioutil"
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
		io.WriteString(rw, "OK!\n")
	}

	oneDone := make(chan struct{})
	one := func(action http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			close(oneDone)
			io.WriteString(rw, "One\n")
			action(rw, req)
		}
	}

	twoDone := make(chan struct{})
	two := func(action http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			close(twoDone)
			io.WriteString(rw, "Two\n")
			action(rw, req)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(NestMiddleware(inner, two, one)))
	defer server.Close()

	res, err := http.Get(server.URL)
	assert.Nil(err)

	<-oneDone
	<-twoDone
	<-innerDone

	contents, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)
	assert.Equal("One\nTwo\nOK!\n", string(contents))
}
