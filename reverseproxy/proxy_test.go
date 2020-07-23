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

var (
	_ webutil.HTTPTracer        = (*mockHTTPTracer)(nil)
	_ webutil.HTTPTraceFinisher = (*mockHTTPTraceFinisher)(nil)
)

func TestProxy(t *testing.T) {
	assert := assert.New(t)

	mockedEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if protoHeader := r.Header.Get(webutil.HeaderXForwardedProto); protoHeader == "" {
			http.Error(w, "No `X-Forwarded-Proto` header!", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ok!")
	}))
	defer mockedEndpoint.Close()

	target, err := url.Parse(mockedEndpoint.URL)
	assert.Nil(err)

	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(target)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)
	assert.Nil(err)

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

func TestProxyTracer(t *testing.T) {
	t.Skip() // these are flaky
	it := assert.New(t)

	target, err := url.Parse("http://web.invalid:9876")
	it.Nil(err)

	tracer := &mockHTTPTracer{}
	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(target)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
		OptProxyTracer(tracer),
	)
	it.Nil(err)
	mockedProxy := httptest.NewServer(proxy)

	res, err := http.Get(mockedProxy.URL)
	it.Nil(err)
	defer res.Body.Close()

	it.Equal(http.StatusBadGateway, res.StatusCode)

	req := tracer.Request
	it.NotNil(req)
	it.Equal("GET", req.Method)
	it.Equal("/", req.URL.String())
	it.Equal(mockedProxy.URL, "http://"+req.Host)

	tf := tracer.Finisher
	it.NotNil(tf)
	it.Nil(tf.Error)
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
		fmt.Fprint(c, "HTTP/1.1 101 Switching Protocols\r\nConnection: upgrade\r\nUpgrade: WebSocket\r\n\r\n")
		bs := bufio.NewScanner(c)
		if !bs.Scan() {
			t.Errorf("backend failed to read line from client: %v", bs.Err())
			return
		}
		fmt.Fprintf(c, "backend got %q\n", bs.Text())
	}))
	defer backendServer.Close()

	backendURL := MustParseURL(backendServer.URL)
	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(backendURL)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)
	assert.Nil(err)

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

	fmt.Fprint(rwc, "Hello\n")
	bs := bufio.NewScanner(rwc)
	assert.True(bs.Scan())

	got := bs.Text()
	want := `backend got "Hello"`
	assert.Equal(got, want)
}

type mockHTTPTracer struct {
	Request  *http.Request
	Finisher *mockHTTPTraceFinisher
}

func (mht *mockHTTPTracer) Start(req *http.Request) (webutil.HTTPTraceFinisher, *http.Request) {
	mht.Request = req
	mht.Finisher = &mockHTTPTraceFinisher{}
	return mht.Finisher, req
}

type mockHTTPTraceFinisher struct {
	Error error
}

func (mhtf *mockHTTPTraceFinisher) Finish(err error) {
	mhtf.Error = err
}
