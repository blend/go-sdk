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

type xmltest struct {
	Foo string `xml:"foo"`
}

func TestXMLResultRender(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	w := webutil.NewMockResponse(buf)
	r := NewCtx(w, webutil.NewMockRequest("GET", "/"))

	xr := &XMLResult{
		StatusCode:	http.StatusOK,
		Response:	xmltest{Foo: "bar"},
	}

	assert.Nil(xr.Render(r))
	assert.Equal(http.StatusOK, w.StatusCode())
	assert.Equal("<xmltest><foo>bar</foo></xmltest>", buf.String())
}

func TestXMLResultRenderStatusCode(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	w := webutil.NewMockResponse(buf)
	r := NewCtx(w, webutil.NewMockRequest("GET", "/"))

	xr := &XMLResult{
		StatusCode:	http.StatusBadRequest,
		Response:	xmltest{Foo: "bar"},
	}

	assert.Nil(xr.Render(r))
	assert.Equal(http.StatusBadRequest, w.StatusCode())
	assert.Equal("<xmltest><foo>bar</foo></xmltest>", buf.String())
}
