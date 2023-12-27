/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"net/http"
	"os"
	"path"
)

// Static returns a static result for an individual file.
func Static(filePath string) *StaticResult {
	file := path.Base(filePath)
	root := path.Dir(filePath)
	return &StaticResult{
		FilePath:	file,
		FileSystem:	http.Dir(root),
	}
}

// StaticResult represents a static output.
type StaticResult struct {
	FilePath	string
	FileSystem	http.FileSystem
	RewriteRules	[]RewriteRule
	Headers		http.Header
}

// Render renders a static result.
func (sr StaticResult) Render(ctx *Ctx) error {
	filePath := sr.FilePath
	for _, rule := range sr.RewriteRules {
		if matched, newFilePath := rule.Apply(filePath); matched {
			filePath = newFilePath
		}
	}

	for key, values := range sr.Headers {
		for _, value := range values {
			ctx.Response.Header().Add(key, value)
		}
	}

	f, err := sr.FileSystem.Open(sr.FilePath)
	if err != nil && os.IsNotExist(err) {
		http.NotFound(ctx.Response, ctx.Request)
		return nil
	} else if err != nil {
		return err
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return err
	}

	http.ServeContent(ctx.Response, ctx.Request, filePath, d.ModTime(), f)
	return nil
}
