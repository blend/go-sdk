/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func handlerNoOp(rw http.ResponseWriter, _ *http.Request, _ *Route, _ RouteParameters) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "OK!\n")
}

func Test_RouteTree_allowed(t *testing.T) {
	its := assert.New(t)

	rt := new(RouteTree)
	rt.Handle(http.MethodGet, "/test", nil)

	allowed := strings.Split(rt.allowed("*", ""), ", ")
	its.Len(allowed, 1)
	its.Equal("GET", allowed[0])

	rt.Handle(http.MethodPost, "/hello", nil)
	allowed = strings.Split(rt.allowed("*", ""), ", ")
	its.Len(allowed, 2)
	its.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == http.MethodGet
	})
	its.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == http.MethodPost
	})

	rt = new(RouteTree)

	rt.Handle(http.MethodGet, "/hello", handlerNoOp)
	allowed = strings.Split(rt.allowed("/hello", ""), ", ")
	its.Len(allowed, 2)
	its.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == "GET"
	})
	its.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == "OPTIONS"
	})
	rt.Handle(http.MethodPost, "/hello", handlerNoOp)
	allowed = strings.Split(rt.allowed("/hello", ""), ", ")
	its.Len(allowed, 3)

	rt.Handle(http.MethodOptions, "/hello", handlerNoOp)
	rt.Handle(http.MethodHead, "/hello", handlerNoOp)
	rt.Handle(http.MethodPut, "/hello", handlerNoOp)
	rt.Handle(http.MethodDelete, "/hello", handlerNoOp)

	rt.Handle(http.MethodPatch, "/hi", handlerNoOp)
	rt.Handle(http.MethodPatch, "/there", handlerNoOp)
	allowed = strings.Split(rt.allowed("/hello", ""), ", ")
	its.Len(allowed, 6)

	rt.Handle(http.MethodPatch, "/hello", handlerNoOp)
	allowed = strings.Split(rt.allowed("/hello", ""), ", ")
	its.Len(allowed, 7)
}

func Test_RouteTree_Route(t *testing.T) {
	its := assert.New(t)

	rt := new(RouteTree)

	rt.Handle(http.MethodGet, "/", handlerNoOp)
	rt.Handle(http.MethodGet, "/foo", handlerNoOp)
	rt.Handle(http.MethodGet, "/foo/:id", handlerNoOp)
	rt.Handle(http.MethodPost, "/foo", handlerNoOp)
	rt.Handle(http.MethodGet, "/bar", handlerNoOp)

	// explicitly register a slash suffixed url here
	rt.Handle(http.MethodGet, "/slash/", handlerNoOp)

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/",
		},
	}
	route, params := rt.Route(req)
	its.NotNil(route)
	its.Equal("/", route.Path)
	its.Empty(params)
	its.Equal("/", req.URL.Path)

	req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/foo",
		},
	}
	route, params = rt.Route(req)
	its.NotNil(route)
	its.Equal("/foo", route.Path)
	its.Equal(http.MethodGet, route.Method)
	its.Empty(params)
	its.Equal("/foo", req.URL.Path)

	req = &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Path: "/foo",
		},
	}
	route, params = rt.Route(req)
	its.NotNil(route)
	its.Equal("/foo", route.Path)
	its.Equal(http.MethodPost, route.Method)
	its.Empty(params)
	its.Equal("/foo", req.URL.Path)

	// explicitly test matching with an extra slash
	req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/foo/",
		},
	}
	route, params = rt.Route(req)
	its.NotNil(route)
	its.Equal("/foo", route.Path)
	its.Empty(params)
	its.Equal("/foo/", req.URL.Path)

	req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/foo/test",
		},
	}
	route, params = rt.Route(req)
	its.NotNil(route)
	its.Equal("/foo/:id", route.Path)
	its.NotEmpty(params)
	its.Equal("test", params["id"])
	its.Equal("/foo/test", req.URL.Path)

	req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/bar",
		},
	}
	route, params = rt.Route(req)
	its.NotNil(route)
	its.Equal("/bar", route.Path)
	its.Empty(params)
	its.Equal("/bar", req.URL.Path)

	req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/slash",
		},
	}
	route, params = rt.Route(req)
	its.NotNil(route)
	its.Equal("/slash/", route.Path)
	its.Empty(params)
	its.Equal("/slash", req.URL.Path)

	req = &http.Request{
		Method: http.MethodConnect,
		URL: &url.URL{
			Path: "/slash",
		},
	}
	route, params = rt.Route(req)
	its.Nil(route)
	its.Empty(params)
	its.Equal("/slash", req.URL.Path)

	req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/slash",
		},
	}
	rt.SkipTrailingSlashRedirects = true
	route, params = rt.Route(req)
	its.Nil(route)
	its.Empty(params)
	its.Equal("/slash", req.URL.Path)
}

func Test_RouteTree_Route_slash(t *testing.T) {
	its := assert.New(t)

	rt := new(RouteTree)

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/",
		},
	}
	route, params := rt.Route(req)
	its.Nil(route)
	its.Empty(params)
	its.Equal("/", req.URL.Path)
}

func Test_RouteTree_withPathAlternateTrailingSlash(t *testing.T) {
	its := assert.New(t)

	its.Equal("/foo", new(RouteTree).withPathAlternateTrailingSlash("/foo/"))
	its.Equal("/foo/", new(RouteTree).withPathAlternateTrailingSlash("/foo"))
	its.Equal("", new(RouteTree).withPathAlternateTrailingSlash(""))
}

func routeExpectsPath(method, path string) Handler {
	return func(rw http.ResponseWriter, req *http.Request, _ *Route, _ RouteParameters) {
		if req.Method != method {
			http.Error(rw, "expects method: "+method, http.StatusBadRequest)
			return
		}
		if req.URL.Path != path {
			http.Error(rw, "expects path: "+path, http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!\n")
	}
}

func callCounter(counter *int32, statusCode int) Handler {
	return func(rw http.ResponseWriter, req *http.Request, _ *Route, _ RouteParameters) {
		defer atomic.AddInt32(counter, 1)
		rw.WriteHeader(statusCode)
		fmt.Fprintf(rw, "counted call!\n")
	}
}

func Test_RouteTree_ServeHTTP(t *testing.T) {
	its := assert.New(t)

	rt := new(RouteTree)

	rt.Handle(http.MethodGet, "/", routeExpectsPath(http.MethodGet, "/"))
	rt.Handle(http.MethodGet, "/foo", routeExpectsPath(http.MethodGet, "/foo"))
	rt.Handle(http.MethodGet, "/foo/:id", routeExpectsPath(http.MethodGet, "/foo/test-id"))
	rt.Handle(http.MethodPost, "/foo", routeExpectsPath(http.MethodPost, "/foo"))
	rt.Handle(http.MethodGet, "/bar", routeExpectsPath(http.MethodGet, "/bar"))

	// explicitly register a slash url here
	rt.Handle(http.MethodGet, "/slash/", handlerNoOp)

	mock := httptest.NewServer(rt)
	defer mock.Close()

	res, err := mock.Client().Get(mock.URL + "/")
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)

	res, err = mock.Client().Get(mock.URL + "/foo")
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)

	res, err = mock.Client().Get(mock.URL + "/foo/")
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)

	res, err = mock.Client().Post(mock.URL+"/foo/", "", nil)
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)

	res, err = mock.Client().Get(mock.URL + "/foo/test-id")
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)

	res, err = mock.Client().Get(mock.URL + "/foo/not-test-id")
	its.Nil(err)
	its.Equal(http.StatusBadRequest, res.StatusCode)

	res, err = mock.Client().Get(mock.URL + "/bar/")
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)

	optionsReq, _ := http.NewRequest(http.MethodOptions, mock.URL, nil)
	// now handle the super weird stuff
	res, err = mock.Client().Do(optionsReq)
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode)
	allowedHeader := res.Header.Get(webutil.HeaderAllow)
	its.NotEmpty(allowedHeader)
	its.Equal("GET, OPTIONS", allowedHeader)

	rt.SkipHandlingMethodOptions = true
	res, err = mock.Client().Do(optionsReq)
	its.Nil(err)
	its.Equal(http.StatusNotFound, res.StatusCode)
	allowedHeader = res.Header.Get(webutil.HeaderAllow)
	its.Empty(allowedHeader)

	var notFoundCalls int32
	rt.NotFoundHandler = callCounter(&notFoundCalls, http.StatusNotFound)
	res, err = mock.Client().Do(optionsReq)
	its.Nil(err)
	its.Equal(http.StatusNotFound, res.StatusCode)
	allowedHeader = res.Header.Get(webutil.HeaderAllow)
	its.Empty(allowedHeader)
	its.Equal(1, notFoundCalls)

	headReq, _ := http.NewRequest(http.MethodHead, mock.URL, nil)
	res, err = mock.Client().Do(headReq)
	its.Nil(err)
	its.Equal(http.StatusMethodNotAllowed, res.StatusCode)
	allowedHeader = res.Header.Get(webutil.HeaderAllow)
	its.NotEmpty(allowedHeader)
	its.Equal("GET, OPTIONS", allowedHeader)

	var methodNotAllowedCalls int32
	rt.MethodNotAllowedHandler = callCounter(&methodNotAllowedCalls, http.StatusMethodNotAllowed)
	res, err = mock.Client().Do(headReq)
	its.Nil(err)
	its.Equal(http.StatusMethodNotAllowed, res.StatusCode)
	allowedHeader = res.Header.Get(webutil.HeaderAllow)
	its.NotEmpty(allowedHeader)
	its.Equal("GET, OPTIONS", allowedHeader)
	its.Equal(1, notFoundCalls)
	its.Equal(1, methodNotAllowedCalls)

	rt.SkipMethodNotAllowed = true
	rt.NotFoundHandler = nil
	res, err = mock.Client().Do(optionsReq)
	its.Nil(err)
	its.Equal(http.StatusNotFound, res.StatusCode)
	allowedHeader = res.Header.Get(webutil.HeaderAllow)
	its.Empty(allowedHeader)

	rt.NotFoundHandler = callCounter(&notFoundCalls, http.StatusNotFound)
	res, err = mock.Client().Do(optionsReq)
	its.Nil(err)
	its.Equal(http.StatusNotFound, res.StatusCode)
	allowedHeader = res.Header.Get(webutil.HeaderAllow)
	its.Empty(allowedHeader)
	its.Equal(2, notFoundCalls)
}
