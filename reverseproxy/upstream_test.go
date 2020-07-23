package reverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestUpstreamWithoutHopHeaders(t *testing.T) {
	assert := assert.New(t)

	u := NewUpstream(MustParseURL("http://localhost:5000"))
	assert.NotNil(u.ReverseProxy)
}

func TestUpstreamServeHTTP(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer srv.Close()

	u := NewUpstream(MustParseURL(srv.URL))
	assert.NotNil(u.ReverseProxy)
	u.Log = logger.None()

	proxy, err := NewProxy()
	assert.Nil(err)
	proxy.Upstreams = append(proxy.Upstreams, u)

	req, err := http.NewRequest("GET", srv.URL, nil)
	assert.Nil(err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)

	assert.Empty(res.Header.Get("X-Forwarded-For"))
	assert.Empty(res.Header.Get("X-Forwarded-Port"))
}
