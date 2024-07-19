/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"net"
	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/webutil"
)

// Raw returns a new raw result.
func Raw(contents []byte) *RawResult {
	return &RawResult{
		StatusCode:  http.StatusOK,
		ContentType: http.DetectContentType(contents),
		Response:    contents,
	}
}

// RawWithContentType returns a binary response with a given content type.
func RawWithContentType(contentType string, body []byte) *RawResult {
	return &RawResult{
		StatusCode:  http.StatusOK,
		ContentType: contentType,
		Response:    body,
	}
}

// RawResult is for when you just want to dump bytes.
type RawResult struct {
	StatusCode  int
	ContentType string
	Response    []byte
}

// Render renders the result.
func (rr *RawResult) Render(ctx *Ctx) error {
	if len(rr.ContentType) != 0 {
		ctx.Response.Header().Set(webutil.HeaderContentType, rr.ContentType)
	}
	ctx.Response.WriteHeader(rr.StatusCode)
	_, err := ctx.Response.Write(rr.Response)
	if err != nil {
		if typed, ok := err.(*net.OpError); ok {
			return ex.New(webutil.ErrNetWrite, ex.OptInner(typed))
		}
		return ex.New(err)
	}
	return nil
}
