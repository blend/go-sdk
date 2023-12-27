/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// NewCachedStaticFile returns a new cached static file for a given path.
func NewCachedStaticFile(path string) (*CachedStaticFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ex.New(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return nil, ex.New(err)
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, ex.New(err)
	}

	return &CachedStaticFile{
		Path:     path,
		Contents: bytes.NewReader(contents),
		ModTime:  finfo.ModTime(),
		ETag:     webutil.ETag(contents),
		Size:     len(contents),
	}, nil
}

var (
	_ Result = (*CachedStaticFile)(nil)
)

// CachedStaticFile is a memory mapped static file.
type CachedStaticFile struct {
	Path     string
	Size     int
	ETag     string
	ModTime  time.Time
	Contents *bytes.Reader
}

// Render implements Result.
//
// Note: It is safe to ingore the error returned from this method; it only
// has this signature to satisfy the `Result` interface.
func (csf CachedStaticFile) Render(ctx *Ctx) error {
	if csf.ETag != "" {
		ctx.Response.Header().Set(webutil.HeaderETag, csf.ETag)
	}
	ctx.WithContext(logger.WithLabel(ctx.Context(), "web.static_file_cached", csf.Path))
	http.ServeContent(ctx.Response, ctx.Request, csf.Path, csf.ModTime, csf.Contents)
	return nil
}
