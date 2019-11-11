package proxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func urlMustParse(urlToParse string) *url.URL {
	url, err := url.Parse(urlToParse)
	if err != nil {
		panic(err)
	}
	return url
}

func TestProxy(t *testing.T) {
	assert := assert.New(t)

	mockedEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if protoHeader := r.Header.Get("X-Forwarded-Proto"); protoHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No `X-Forwarded-Proto` header!"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok!"))
		return
	}))
	defer mockedEndpoint.Close()

	target, err := url.Parse(mockedEndpoint.URL)
	assert.Nil(err)

	proxy := New().WithUpstream(NewUpstream(target))
	proxy.WithUpstreamHeader("X-Forwarded-Proto", "http")

	mockedProxy := httptest.NewServer(proxy)

	res, err := http.Get(mockedProxy.URL)
	assert.Nil(err)
	defer res.Body.Close()

	fullBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	mockedContents := string(fullBody)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("Ok!", mockedContents)
}
