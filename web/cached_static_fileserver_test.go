package web

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestCachedStaticFileserver(t *testing.T) {
	assert := assert.New(t)

	cfs := NewCachedStaticFileServer(http.Dir("testdata"))
	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req).WithRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	}))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes())

	assert.NotEmpty(cfs.files)
	file, hasFile := cfs.files["test_file.html"]
	assert.True(hasFile)
	assert.NotNil(file)
}

func TestCachedStaticFileserverHeaders(t *testing.T) {
	assert := assert.New(t)

	cfs := NewCachedStaticFileServer(http.Dir("testdata"))
	cfs.AddHeader("foo", "bar")
	assert.NotEmpty(cfs.Headers())

	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req).WithRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	}))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes())
	assert.Equal("bar", res.Header().Get("foo"), "the header should be set on the response")
	assert.NotEmpty(res.Header().Get(HeaderETag))
}

func TestCachedStaticFileserverRewriteRule(t *testing.T) {
	assert := assert.New(t)

	cfs := NewCachedStaticFileServer(http.Dir("testdata"))
	assert.Nil(cfs.AddRewriteRule(RegexpAssetCacheFiles, func(path string, parts ...string) string {
		return fmt.Sprintf("%s.%s", parts[1], parts[3])
	}))

	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.123123123.html")
	result := cfs.Action(NewCtx(res, req).WithRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.123123123.html",
	}))

	assert.Nil(result)
	assert.NotEmpty(buffer.Bytes(), "we should still have reached the file")
}

func TestCachedStaticFileserverMiddleware(t *testing.T) {
	assert := assert.New(t)

	var didCallMiddleware bool
	var didNestMiddleware bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	cfs := NewCachedStaticFileServer(http.Dir("testdata"))
	cfs.SetMiddleware(func(action Action) Action {
		didNestMiddleware = true
		return func(ctx *Ctx) Result {
			defer wg.Done()
			didCallMiddleware = true
			return action(ctx)
		}
	})

	buffer := bytes.NewBuffer(nil)
	res := webutil.NewMockResponse(buffer)
	req := webutil.NewMockRequest("GET", "/test_file.html")
	result := cfs.Action(NewCtx(res, req).WithRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	}))
	wg.Wait()

	assert.Nil(result)
	assert.True(didNestMiddleware)
	assert.True(didCallMiddleware)
	assert.NotEmpty(buffer.Bytes())

	didNestMiddleware = false
	didCallMiddleware = false
	wg.Add(1)
	result = cfs.Action(NewCtx(res, req).WithRouteParams(RouteParameters{
		RouteTokenFilepath: "test_file.html",
	}))
	wg.Wait()

	assert.Nil(result)
	assert.False(didNestMiddleware)
	assert.True(didCallMiddleware)
	assert.NotEmpty(buffer.Bytes())
}
