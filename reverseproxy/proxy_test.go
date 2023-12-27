/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reverseproxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

var (
	_	webutil.HTTPTracer		= (*mockHTTPTracer)(nil)
	_	webutil.HTTPTraceFinisher	= (*mockHTTPTraceFinisher)(nil)
)

func Test_Proxy(t *testing.T) {
	its := assert.New(t)

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
	its.Nil(err)

	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(target)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)
	its.Nil(err)

	mockedProxy := httptest.NewServer(proxy)
	defer mockedProxy.Close()

	res, err := http.Get(mockedProxy.URL)
	its.Nil(err)
	defer res.Body.Close()

	its.Empty(res.Header.Get("x-forwarded-proto"))
	its.Empty(res.Header.Get("x-forwarded-port"))

	fullBody, err := io.ReadAll(res.Body)
	its.Nil(err)

	mockedContents := string(fullBody)
	its.Equal(http.StatusOK, res.StatusCode)
	its.Equal("Ok!", mockedContents)
}

func Test_Proxy_Tracer(t *testing.T) {
	its := assert.New(t)

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
	its.Nil(err)

	tracer := &mockHTTPTracer{}
	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(target)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
		OptProxyTracer(tracer),
	)
	its.Nil(err)

	mockedProxy := httptest.NewServer(proxy)
	defer mockedProxy.Close()

	res, err := http.Get(mockedProxy.URL)
	its.Nil(err)
	defer res.Body.Close()

	its.Equal(http.StatusOK, res.StatusCode)

	req := tracer.Request
	its.NotNil(req)
	its.Equal("GET", req.Method)
	its.Equal("/", req.URL.String())
	its.Equal(mockedProxy.URL, "http://"+req.Host)

	its.Equal(http.StatusOK, tracer.StatusCode)
	its.Nil(tracer.Error)
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
	Request		*http.Request
	StatusCode	int
	Error		error
}

func (mht *mockHTTPTracer) Start(req *http.Request) (webutil.HTTPTraceFinisher, *http.Request) {
	mht.Request = req
	return &mockHTTPTraceFinisher{mht}, req
}

type mockHTTPTraceFinisher struct {
	Tracer *mockHTTPTracer
}

func (mhtf *mockHTTPTraceFinisher) Finish(statusCode int, err error) {
	mhtf.Tracer.StatusCode = statusCode
	mhtf.Tracer.Error = err
}

func TestProxy_Panic(t *testing.T) {
	its := assert.New(t)

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
	its.Nil(err)

	log := logger.Memory(io.Discard)
	defer log.Close()

	errors := make(chan error)
	log.Listen(logger.Fatal, "panic-chan", logger.NewErrorEventListener(func(ctx context.Context, e logger.ErrorEvent) {
		errors <- e.Err
	}))

	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(
			target,
		)),
		OptProxyLog(log),
		OptProxyResolver(func(_ *http.Request, _ []*Upstream) (*Upstream, error) {
			panic("this is just a test")
		}),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)
	its.Nil(err)

	mockedProxy := httptest.NewServer(proxy)

	res, err := http.Get(mockedProxy.URL)
	its.Nil(err)
	defer res.Body.Close()
	its.Equal(http.StatusOK, res.StatusCode)
	err = <-errors
	its.NotNil(err)
	its.Equal("this is just a test", err.Error())
}

func TestProxy_Panic_httpAbortHandler(t *testing.T) {
	its := assert.New(t)

	var didCallEndpoint bool
	mockedEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { didCallEndpoint = true }()
		if protoHeader := r.Header.Get(webutil.HeaderXForwardedProto); protoHeader == "" {
			http.Error(w, "No `X-Forwarded-Proto` header!", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ok!")
	}))
	defer mockedEndpoint.Close()

	target, err := url.Parse(mockedEndpoint.URL)
	its.Nil(err)

	log := logger.Memory(io.Discard)
	defer log.Close()

	errors := make(chan error, 1)
	log.Listen(logger.Fatal, "panic-chan", logger.NewErrorEventListener(func(ctx context.Context, e logger.ErrorEvent) {
		errors <- e.Err
	}))

	proxy, err := NewProxy(
		OptProxyUpstream(NewUpstream(
			target,
		)),
		OptProxyLog(log),
		OptProxyResolver(func(_ *http.Request, _ []*Upstream) (*Upstream, error) {
			panic(http.ErrAbortHandler)
		}),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
	)
	its.Nil(err)

	mockedProxy := httptest.NewServer(proxy)

	res, err := http.Get(mockedProxy.URL)
	its.Nil(err)
	defer res.Body.Close()
	its.Equal(http.StatusOK, res.StatusCode)

	// explicitly drain so we process any errors that would come up
	log.Drain()

	its.Empty(errors)
	its.False(didCallEndpoint)
}
