/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/ex"
)

// KnownExtenions are known extenions mapped to their content types.
var (
	KnownExtensions = map[string]string{
		".html": "text/html; charset=utf-8",
		".xml":  "text/xml; charset",
		".json": "application/json; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
	}
)

// DetectContentType generates the content type of a given file by path.
func DetectContentType(path string) (string, error) {
	if contentType, ok := KnownExtensions[filepath.Ext(path)]; ok {
		return contentType, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return "", ex.New(err)
	}
	defer f.Close()
	header := make([]byte, 512)
	_, err = f.Read(header)
	if err != nil {
		return "", ex.New(err)
	}
	return http.DetectContentType(header), nil
}
