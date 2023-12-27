/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
	)
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

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
		OptStaticFileServerHeaders(http.Header{"buzz": []string{"fuzz"}}),
	)
	cfs.AddHeader("foo", "bar")
	assert.NotEmpty(cfs.Headers)

	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	})))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes())

	assert.Equal("bar", res.Header().Get("foo"), "the header should be set on the response")
	assert.Equal("fuzz", res.Header().Get("buzz"), "the header should be set on the response")
}

func TestStaticFileserverRewriteRule(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
	)
	assert.Nil(cfs.AddRewriteRule(RegexpAssetCacheFiles, func(path string, parts ...string) string {
		return fmt.Sprintf("%s.%s", parts[1], parts[3])
	}))

	buffer := new(bytes.Buffer)
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

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
	)
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/"+uuid.V4().String())
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	})))

	assert.Nil(result)
	assert.Equal(http.StatusNotFound, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
}

func TestStaticFileserverNotFoundDefaultProvider(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
	)
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/"+uuid.V4().String())
	result := cfs.Action(NewCtx(res, req,
		OptCtxRouteParams(RouteParameters{RouteTokenFilepath: req.URL.Path}),
		OptCtxDefaultProvider(JSON),
	))

	assert.NotNil(result)
	typed, ok := result.(*JSONResult)
	assert.True(ok)
	assert.NotNil(typed)
	assert.Equal(http.StatusNotFound, typed.StatusCode)
}

func TestStaticFileserverLive(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
	)
	cfs.CacheDisabled = true
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	})))

	assert.Nil(result)
	assert.Equal(http.StatusOK, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
}

func TestStaticFileserverLiveNotFound(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
		OptStaticFileServerCacheDisabled(true),
	)
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/"+uuid.V4().String())
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	})))

	assert.Nil(result)
	assert.Equal(http.StatusNotFound, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
}

func TestStaticFileserverCachedNotFoundOnRoot(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
		OptStaticFileServerCacheDisabled(false),
	)
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/")
	r := NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	}))
	r.DefaultProvider = Text
	result := cfs.Action(r)

	assert.NotNil(result)
	assert.Nil(result.Render(r))
	assert.Equal(http.StatusNotFound, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
}

func TestStaticFileserverLiveNotFoundOnRoot(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
		OptStaticFileServerCacheDisabled(true),
	)
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/")
	r := NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	}))
	r.DefaultProvider = Text
	result := cfs.Action(r)

	assert.NotNil(result)
	assert.Nil(result.Render(r))
	assert.Equal(http.StatusNotFound, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
}

func TestStaticFileserverAddsETag(t *testing.T) {
	assert := assert.New(t)

	cfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(http.Dir("testdata")),
		OptStaticFileServerCacheDisabled(false),
	)
	buffer := new(bytes.Buffer)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req, OptCtxRouteParams(RouteParameters{
		RouteTokenFilepath: req.URL.Path,
	})))

	assert.Nil(result)
	assert.Equal(http.StatusOK, res.StatusCode())
	assert.NotEmpty(buffer.Bytes())
	assert.NotEmpty(res.Header().Get(webutil.HeaderETag))
}
