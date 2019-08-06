package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
)

// RequestFactory creates a new request with a given set of options.
type RequestFactory []r2.Option

// New creates a new request.
func (rf RequestFactory) New(target string, options ...r2.Option) *r2.Request {
	return r2.New(target, append(rf, options...)...)
}

func (rf RequestFactory) Google() (*http.Response, error) {
	return rf.New("https://google.com/robots.txt",
		r2.OptUserAgent("blend go-sdk"),
		r2.OptTimeout(5*time.Second),
	).DiscardWithResponse()
}

func TestMockedRequest(t *testing.T) {
	assert := assert.New(t)

	var didCallHandler bool
	mockServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		didCallHandler = true
		rw.WriteHeader(200)
		fmt.Fprint(rw, "OK!\n")
		return
	}))
	defer mockServer.Close()

	rf := RequestFactory([]r2.Option{r2.OptURL(mockServer.URL)})

	res, err := rf.Google()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.True(didCallHandler)
}
