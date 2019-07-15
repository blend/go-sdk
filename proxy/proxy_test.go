package proxy

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/net/http/httpguts"
	"io"
	"io/ioutil"
	"net"
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
	if !httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
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
	assert := assert.New(t)

	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(upgradeType(r.Header), "websocket")

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

	backendURL := urlMustParse(backendServer.URL)
	proxy := New().WithUpstream(NewUpstream(backendURL))
	proxy.WithUpstreamHeader(webutil.HeaderXForwardedProto, webutil.SchemeHTTP)

	frontendProxy := httptest.NewServer(proxy)
	defer frontendProxy.Close()

	req, _ := http.NewRequest("GET", frontendProxy.URL, nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")

	c := frontendProxy.Client()
	res, err := c.Do(req)
	assert.Nil(err)

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

// Referencing https://github.com/samsarahq/oauth2_proxy/blob/master/websocket_test.go#L13
func TestReverseProxyWebsocketWithScheme(t *testing.T) {
	assert := assert.New(t)

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
			return
		}

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
			return
		}

		if err = conn.WriteMessage(messageType, p); err != nil {
			t.Error(err)
			return
		}
	}))
	defer backend.Close()

	backendURL := urlMustParse(backend.URL)
	backendHostname, backendPort, _ := net.SplitHostPort(backendURL.Host)
	backendHost := net.JoinHostPort(backendHostname, backendPort)
	proxyURL := urlMustParse(backendURL.Scheme + "://" + backendHost + "/")

	proxy := New().WithUpstream(NewUpstream(proxyURL))
	frontend := httptest.NewServer(proxy)
	defer frontend.Close()

	wsUrl := new(url.URL)
	*wsUrl = *proxyURL
	wsUrl.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	assert.Nil(err)

	expected := "hello world"
	err = conn.WriteMessage(websocket.TextMessage, []byte(expected))
	assert.Nil(err)

	messageType, received, err := conn.ReadMessage()
	assert.Nil(err)
	assert.Equal(messageType, websocket.TextMessage)
	assert.Equal(string(received), expected)
}
