package proxy

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func urlMustParse(urlToParse string) *url.URL {
	url, err := url.Parse(urlToParse)
	if err != nil {
		panic(err)
	}
	return url
}

func upgradeType(h http.Header) string {
	if value := h.Get(http.CanonicalHeaderKey("Connection")); value != "" {
		if value != "Upgrade" {
			return ""
		}
	}
	return strings.ToLower(h.Get("Upgrade"))
}

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

	proxy := New().WithUpstream(NewUpstream(target))
	proxy.WithUpstreamHeader(webutil.HeaderXForwardedProto, webutil.SchemeHTTP)

	mockedProxy := httptest.NewServer(proxy)
	defer mockedProxy.Close()

	res, err := http.Get(mockedProxy.URL)
	assert.Nil(err)
	defer res.Body.Close()

	fullBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	mockedContents := string(fullBody)
	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("Ok!", mockedContents)
}

// Referencing https://golang.org/src/net/http/httputil/reverseproxy_test.go
func TestReverseProxyWebSocket(t *testing.T) {
	t.Skip()
	assert := assert.New(t)

	var upgradeTypeValue string
	var hijackErr error
	var scanErr error
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request headers: %#v\n", r.Header)
		upgradeTypeValue = upgradeType(r.Header)
		c, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			hijackErr = err
			return
		}
		defer c.Close()
		io.WriteString(c, "HTTP/1.1 101 Switching Protocols\r\nConnection: upgrade\r\nUpgrade: WebSocket\r\n\r\n")
		bs := bufio.NewScanner(c)
		if !bs.Scan() {
			scanErr = fmt.Errorf("backend failed to read line from client: %v", bs.Err())
			return
		}
		fmt.Fprintf(c, "backend got %q\n", bs.Text())
	}))
	defer backendServer.Close()

	backendURL := urlMustParse(backendServer.URL)
	proxy := New().WithUpstream(NewUpstream(backendURL))
	proxy.WithUpstreamHeader(webutil.HeaderXForwardedProto, webutil.SchemeHTTP)

	frontendProxy := httptest.NewServer(proxy)
	defer frontendProxy.Close()

	req, err := http.NewRequest("GET", frontendProxy.URL, nil)
	assert.Nil(err)
	req.Header = http.Header{
		http.CanonicalHeaderKey("Connection"): []string{"Upgrade"},
		http.CanonicalHeaderKey("Upgrade"):    []string{"websocket"},
	}
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)

	assert.Equal("upgrade", upgradeTypeValue)
	assert.Nil(hijackErr)
	assert.Nil(scanErr)

	assert.Equal(res.StatusCode, 101)
	assert.Equal(upgradeType(req.Header), "websocket")

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
