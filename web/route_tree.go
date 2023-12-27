/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// RouteTree implements the basic logic of creating a route tree.
//
// It is embedded in *App and handles the path to handler matching based
// on route trees per method.
//
// A very simple example:
//
//	rt := new(web.RouteTree)
//	rt.Handle(http.MethodGet, "/", func(w http.ResponseWriter, req *http.Request, route *web.Route, params web.Params) {
//	    w.WriteHeader(http.StatusOK)
//	    fmt.Fprintf(w, "OK!")
//	})
//	(&http.Server{Addr: "127.0.0.1:8080", Handler: rt}).ListenAndServe()
type RouteTree struct {
	// Routes is a map between canonicalized http method
	// (i.e. `GET` vs. `get`) and individual method
	// route trees.
	Routes	map[string]*RouteNode
	// SkipTrailingSlashRedirects disables matching
	// routes that are off by a trailing slash, either because
	// routes are registered with the '/' suffix, or because
	// the request has a '/' suffix and the
	// registered route does not.
	SkipTrailingSlashRedirects	bool
	// SkipHandlingMethodOptions disables returning
	// a result with the `ALLOWED` header for method options,
	// and will instead 404 for `OPTIONS` methods.
	SkipHandlingMethodOptions	bool
	// SkipMethodNotAllowed skips specific handling
	// for methods that do not have a route tree with
	// a specific 405 response, and will instead return a 404.
	SkipMethodNotAllowed	bool
	// NotFoundHandler is an optional handler to set
	// to customize not found (404) results.
	NotFoundHandler	Handler
	// MethodNotAllowedHandler is an optional handler
	// to set to customize method not allowed (405) results.
	MethodNotAllowedHandler	Handler
}

// Handle adds a handler at a given method and path.
func (rt *RouteTree) Handle(method, path string, handler Handler) {
	if len(path) == 0 {
		panic("path must not be empty")
	}
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if rt.Routes == nil {
		rt.Routes = make(map[string]*RouteNode)
	}

	root := rt.Routes[method]
	if root == nil {
		root = new(RouteNode)
		rt.Routes[method] = root
	}
	root.AddRoute(method, path, handler)
}

// Route gets the route and parameters for a given request
// if it matches a registered handler.
//
// It will automatically resolve if a trailing slash should be appended
// for the input request url path, and will return the corresponding redirected
// route (and parameters) if there is one.
func (rt *RouteTree) Route(req *http.Request) (*Route, RouteParameters) {
	path := req.URL.Path
	methodRoot := rt.Routes[req.Method]
	if methodRoot != nil {
		route, params, shouldRedirectTrailingSlash := methodRoot.getValue(path)
		if req.Method != http.MethodConnect && path != "/" {
			if shouldRedirectTrailingSlash && !rt.SkipTrailingSlashRedirects {
				route, params, _ = methodRoot.getValue(rt.withPathAlternateTrailingSlash(path))
			}
		}
		return route, params
	}
	return nil, nil
}

// ServeHTTP makes the router implement the http.Handler interface.
func (rt *RouteTree) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if root := rt.Routes[req.Method]; root != nil {
		route, params, trailingSlashRedirect := root.getValue(path)
		if route != nil {
			route.Handler(w, req, route, params)
			return
		} else if req.Method != http.MethodConnect && path != "/" {
			if trailingSlashRedirect && !rt.SkipTrailingSlashRedirects {
				rt.redirectTrailingSlash(w, req)
				return
			}
		}
	}

	if req.Method == http.MethodOptions {
		// Handle OPTIONS requests
		if !rt.SkipHandlingMethodOptions {
			if allow := rt.allowed(path, req.Method); allow != "" {
				w.Header().Set(webutil.HeaderAllow, allow)
				// just return the allowed header
				return
			}
			// return a 404 below
		}
	} else {
		// Handle 405
		if !rt.SkipMethodNotAllowed {
			if allow := rt.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set(webutil.HeaderAllow, allow)
				if rt.MethodNotAllowedHandler != nil {
					rt.MethodNotAllowedHandler(w, req, nil, nil)
					return
				}
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
		}
	}

	// Handle 404
	if rt.NotFoundHandler != nil {
		rt.NotFoundHandler(w, req, nil, nil)
	} else {
		http.NotFound(w, req)
	}
}

//
// internal helpers
//

// withPathAlternateTrailingSlash returns the request with a `/` suffix on the url path.
func (rt *RouteTree) withPathAlternateTrailingSlash(path string) string {
	// if the path has a slash already, try removing it
	if len(path) > 1 && path[len(path)-1] == '/' {
		// try removing the slash
		return path[:len(path)-1]
	}
	if len(path) > 0 {
		// try adding the slash
		return path + "/"
	}
	return path
}

// redirectTrailingSlash redirects the request if a suffix trailing
// forward slash should be added.
func (rt *RouteTree) redirectTrailingSlash(w http.ResponseWriter, req *http.Request) {
	code := http.StatusMovedPermanently	// 301 // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect	// 307
	}
	req.URL.Path = rt.withPathAlternateTrailingSlash(req.URL.Path)
	http.Redirect(w, req, req.URL.String(), code)
	return
}

func (rt *RouteTree) allowed(path, reqMethod string) (allow string) {
	if path == "*" {	// server-wide
		for method := range rt.Routes {
			if method == http.MethodOptions {
				continue
			}

			// add request method to list of allowed methods
			if allow == "" {
				allow = method
			} else {
				allow += ", " + method
			}
		}
		return
	}
	for method := range rt.Routes {
		// Skip the requested method - we already tried this one
		if method == reqMethod || method == http.MethodOptions {
			continue
		}

		handle, _, _ := rt.Routes[method].getValue(path)
		if handle != nil {
			// add request method to list of allowed methods
			if allow == "" {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	}
	if allow != "" {
		allow += ", " + http.MethodOptions
	}
	return
}
