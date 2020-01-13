package reverseproxy

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestProxy(t *testing.T) {
	assert := assert.New(t)

	mockedEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if protoHeader := r.Header.Get(webutil.HeaderXForwardedProto); protoHeader == "" {
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

	proxy := NewProxy(
		OptProxyUpstream(NewUpstream(target)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)

	mockedProxy := httptest.NewServer(proxy)

	res, err := http.Get(mockedProxy.URL)
	assert.Nil(err)
	defer res.Body.Close()

	assert.Empty(res.Header.Get("x-forwarded-proto"))
	assert.Empty(res.Header.Get("x-forwarded-port"))

	fullBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	mockedContents := string(fullBody)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("Ok!", mockedContents)
}

// Referencing https://golang.org/src/net/http/httputil/reverseproxy_test.go
func TestReverseProxyWebSocket(t *testing.T) {
	assert := assert.New(t)

	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(UpgradeType(r.Header), "websocket")

		c, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
		io.WriteString(c, "HTTP/1.1 101 Switching Protocols\r\nConnection: upgrade\r\nUpgrade: WebSocket\r\n\r\n")
		bs := bufio.NewScanner(c)
		if !bs.Scan() {
			t.Errorf("backend failed to read line from client: %v", bs.Err())
			return
		}
		fmt.Fprintf(c, "backend got %q\n", bs.Text())
	}))
	defer backendServer.Close()

	backendURL := MustParseURL(backendServer.URL)
	proxy := NewProxy(
		OptProxyUpstream(NewUpstream(backendURL)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)

	frontendProxy := httptest.NewServer(proxy)
	defer frontendProxy.Close()

	req, _ := http.NewRequest("GET", frontendProxy.URL, nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")

	c := frontendProxy.Client()
	res, err := c.Do(req)
	assert.Nil(err)

	assert.Equal(res.StatusCode, 101)

	assert.Equal(UpgradeType(req.Header), "websocket")

	rwc, ok := res.Body.(io.ReadWriteCloser)
	assert.True(ok)
	defer rwc.Close()

	io.WriteString(rwc, "Hello\n")
	bs := bufio.NewScanner(rwc)
	assert.True(bs.Scan())

	got := bs.Text()
	want := `backend got "Hello"`
	assert.Equal(got, want)
}
