/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/logger"
)

// Redirect returns a redirect result to a given destination.
func Redirect(destination string) *RedirectResult {
	return &RedirectResult{
		RedirectURI: destination,
	}
}

// Redirectf returns a redirect result to a given destination specified by a given format and scan arguments.
func Redirectf(format string, args ...interface{}) *RedirectResult {
	return &RedirectResult{
		RedirectURI: fmt.Sprintf(format, args...),
	}
}

// RedirectWithMethod returns a redirect result to a destination with a given method.
func RedirectWithMethod(method, destination string) *RedirectResult {
	return &RedirectResult{
		Method:      method,
		RedirectURI: destination,
	}
}

// RedirectWithMethodf returns a redirect result to a destination composed of a format and scan arguments with a given method.
func RedirectWithMethodf(method, format string, args ...interface{}) *RedirectResult {
	return &RedirectResult{
		Method:      method,
		RedirectURI: fmt.Sprintf(format, args...),
	}
}

// RedirectResult is a result that should cause the browser to redirect.
type RedirectResult struct {
	Method      string `json:"redirect_method"`
	RedirectURI string `json:"redirect_uri"`
}

// Render writes the result to the response.
func (rr *RedirectResult) Render(ctx *Ctx) error {
	ctx.WithContext(logger.WithLabel(ctx.Context(), "web.redirect", rr.RedirectURI))
	if len(rr.Method) > 0 {
		ctx.Request.Method = rr.Method
		http.Redirect(ctx.Response, ctx.Request, rr.RedirectURI, http.StatusFound)
	} else {
		http.Redirect(ctx.Response, ctx.Request, rr.RedirectURI, http.StatusTemporaryRedirect)
	}
	return nil
}
