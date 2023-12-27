/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestViewModelWrap(t *testing.T) {
	assert := assert.New(t)

	indexTemplate := `{{ define "index" }}{{ range $index, $obj := .ViewModel }}<div>{{ template "control" ( $.Wrap $obj ) }}</div>{{ end }}{{ end }}`
	controlTemplate := `{{ define "control" }}{{ if .Ctx }}{{ .ViewModel }}{{ end }}{{ end }}`

	app := MustNew()
	app.Views.AddLiterals(indexTemplate, controlTemplate)

	app.GET("/", func(r *Ctx) Result {
		return r.Views.View("index", []string{"foo", "bar", "baz"})
	})

	contents, meta, err := MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Equal("<div>foo</div><div>bar</div><div>baz</div>", string(contents))
}
