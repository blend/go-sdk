/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func Test_NewCachedStaticFile(t *testing.T) {
	its := assert.New(t)

	csf, err := NewCachedStaticFile("testdata/test_file.html")
	its.Nil(err)
	its.Equal("testdata/test_file.html", csf.Path)
	its.Equal(190, csf.Size)
	its.False(csf.ModTime.IsZero())
	its.Equal("da9a836ffc32feea4b26a536d3d0eccc", csf.ETag)

	contents, err := io.ReadAll(csf.Contents)
	its.Nil(err)
	its.Contains(string(contents), `<title>Test!</title>`)
}

func Test_CachedStaticFile_Render(t *testing.T) {
	its := assert.New(t)

	csf, err := NewCachedStaticFile("testdata/test_file.html")
	its.Nil(err)

	buf := new(bytes.Buffer)
	ctx := MockCtxWithBuffer(http.MethodGet, "index.html", buf)
	err = csf.Render(ctx)
	its.Nil(err)
	its.Equal("da9a836ffc32feea4b26a536d3d0eccc", ctx.Response.Header().Get(webutil.HeaderETag))
	its.Equal("text/html; charset=utf-8", ctx.Response.Header().Get(webutil.HeaderContentType))

	its.Equal("testdata/test_file.html", logger.GetLabels(ctx.Context())["web.static_file_cached"])
	its.Contains(buf.String(), `<title>Test!</title>`)
}
