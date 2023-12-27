/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestJSONResultRender(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	w := webutil.NewMockResponse(buf)
	r := NewCtx(w, webutil.NewMockRequest("GET", "/"))

	jr := &JSONResult{
		StatusCode: http.StatusOK,
		Response: map[string]interface{}{
			"foo": "bar",
		},
	}

	assert.Nil(jr.Render(r))
	assert.Equal(http.StatusOK, w.StatusCode())
	assert.Equal("{\"foo\":\"bar\"}\n", buf.String())
}

func TestJSONResultRenderStatusCode(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	w := webutil.NewMockResponse(buf)
	r := NewCtx(w, webutil.NewMockRequest("GET", "/"))

	jr := &JSONResult{
		StatusCode: http.StatusBadRequest,
		Response: map[string]interface{}{
			"foo": "bar",
		},
	}

	assert.Nil(jr.Render(r))
	assert.Equal(http.StatusBadRequest, w.StatusCode())
	assert.Equal("{\"foo\":\"bar\"}\n", buf.String())
}
