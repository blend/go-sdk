/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestRedirectResult(t *testing.T) {
	assert := assert.New(t)

	resBody := new(bytes.Buffer)
	res := webutil.NewMockResponse(resBody)
	req := webutil.NewMockRequest("GET", "/")
	ctx := NewCtx(res, req)

	assert.Nil((&RedirectResult{RedirectURI: "/foo"}).Render(ctx))

	assert.Equal(http.StatusTemporaryRedirect, res.StatusCode())
	assert.Contains(resBody.String(), "/foo", resBody.String())
}

func TestRedirectResultMethod(t *testing.T) {
	assert := assert.New(t)

	resBody := new(bytes.Buffer)
	res := webutil.NewMockResponse(resBody)
	req := webutil.NewMockRequest("POST", "/")
	ctx := NewCtx(res, req)

	assert.Nil((&RedirectResult{Method: "GET", RedirectURI: "/foo"}).Render(ctx))

	assert.Equal(http.StatusFound, res.StatusCode())
	assert.Contains(resBody.String(), "/foo", resBody.String())
}
