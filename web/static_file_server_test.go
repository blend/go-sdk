package web

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
)

func TestStaticFileserver(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(http.Dir("testdata"))
	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	})))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes())
}

func TestStaticFileserverHeaders(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(http.Dir("testdata"))
	cfs.AddHeader("foo", "bar")
	assert.NotEmpty(cfs.Headers)

	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	})))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes())

	assert.Equal("bar", res.Header().Get("foo"), "the header should be set on the response")
}

func TestStaticFileserverRewriteRule(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(http.Dir("testdata"))
	assert.Nil(cfs.AddRewriteRule(RegexpAssetCacheFiles, func(path string, parts ...string) string {
		return fmt.Sprintf("%s.%s", parts[1], parts[3])
	}))

	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.123123123.html")
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.123123123.html",
	})))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes(), "we should still have reached the file")
}

func TestStaticFileserverNotFound(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(http.Dir("testdata"))
	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/"+uuid.V4().String())
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	})))

	assert.Nil(result)
	assert.Equal(http.StatusNotFound, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
}
