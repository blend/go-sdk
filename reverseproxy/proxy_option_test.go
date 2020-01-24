package reverseproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestOptProxyTransformRequest(t *testing.T) {
	it := assert.New(t)

	var requests []*http.Request
	tr := func(req *http.Request) {
		requests = append(requests, req)
	}

	target, err := url.Parse("http://web.invalid:9876")
	it.Nil(err)

	p := NewProxy(
		OptProxyUpstream(NewUpstream(target)),
		OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTP),
		OptProxyTransformRequest(tr),
	)
	// Need to special case function equality.
	it.Equal(reflect.ValueOf(tr).Pointer(), reflect.ValueOf(p.TransformRequest).Pointer())

	mockedProxy := httptest.NewServer(p)
	res, err := http.Get(mockedProxy.URL)
	it.Nil(err)
	defer res.Body.Close()
	it.Equal(http.StatusBadGateway, res.StatusCode)

	it.Len(requests, 1)
	calledReq := requests[0]
	it.NotNil(calledReq)
	it.Equal("GET", calledReq.Method)
	it.Equal("/", calledReq.URL.String())
	it.Equal(mockedProxy.URL, "http://"+calledReq.Host)
}
