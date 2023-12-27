/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestViewCacheProperties(t *testing.T) {
	assert := assert.New(t)

	vc, err := NewViewCache()
	assert.Nil(err)
	assert.NotNil(vc.FuncMap)

	assert.Equal(DefaultTemplateNameBadRequest, vc.BadRequestTemplateName)
	assert.Nil(OptViewCacheBadRequestTemplateName("foo")(vc))
	assert.Equal("foo", vc.BadRequestTemplateName)

	assert.Equal(DefaultTemplateNameInternalError, vc.InternalErrorTemplateName)
	assert.Nil(OptViewCacheInternalErrorTemplateName("bar")(vc))
	assert.Equal("bar", vc.InternalErrorTemplateName)

	assert.Equal(DefaultTemplateNameNotFound, vc.NotFoundTemplateName)
	assert.Nil(OptViewCacheNotFoundTemplateName("baz")(vc))
	assert.Equal("baz", vc.NotFoundTemplateName)

	assert.Equal(DefaultTemplateNameNotAuthorized, vc.NotAuthorizedTemplateName)
	assert.Nil(OptViewCacheNotAuthorizedTemplateName("buzz")(vc))
	assert.Equal("buzz", vc.NotAuthorizedTemplateName)

	assert.Equal(DefaultTemplateNameStatus, vc.StatusTemplateName)
	assert.Nil(OptViewCacheStatusTemplateName("fuzz")(vc))
	assert.Equal("fuzz", vc.StatusTemplateName)

	assert.Empty(vc.Paths)
	assert.Nil(OptViewCachePaths("foo", "bar")(vc))
	assert.NotEmpty(vc.Paths)

	assert.Empty(vc.Literals)
	assert.Nil(OptViewCacheLiterals("boo", "loo")(vc))
	assert.NotEmpty(vc.Literals)
}

func TestViewCacheAddRawViews(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	vc.AddLiterals(`{{ define "test" }}<h1> This is a test. </h1>{{ end }}`)

	view, err := vc.Parse()
	assert.Nil(err)
	assert.NotNil(view)

	buf := bytes.NewBuffer(nil)
	assert.Nil(view.ExecuteTemplate(buf, "test", nil))
	assert.NotEmpty(buf.String())

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	_, err = vc.Parse()
	assert.NotNil(err)

	vc = MustNewViewCache()
	vc.AddPaths("this path doesn't exist at all")
	_, err = vc.Parse()
	assert.NotNil(err)
}

func TestViewCacheCached(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.False(vc.LiveReload)
	vc.AddLiterals(`{{ define "foo" }}bar{{ end }}`)
	assert.Nil(vc.Initialize())

	tmp, err := vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf := bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc.Literals = []string{`{{ define "foo" }}baz{{ end }}`}
	tmp, err = vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf = bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	_, err = vc.Lookup("foo")
	assert.NotNil(err)
}

func TestViewCacheCachingDisabled(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache(OptViewCacheLiveReload(true))
	assert.True(vc.LiveReload)
	vc.AddLiterals(`{{ define "foo" }}bar{{ end }}`)
	assert.Nil(vc.Initialize())

	tmp, err := vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf := bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc.Literals = []string{`{{ define "foo" }}baz{{ end }}`}
	tmp, err = vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf = bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("baz", buf.String())
}

func TestViewCacheBadRequest(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.Nil(vc.Initialize())
	vr, _ := vc.BadRequest(ex.Class("only a test")).(*ViewResult)
	assert.Equal(vc.BadRequestTemplateName, vr.ViewName)
	assert.Equal(http.StatusBadRequest, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.ViewModel)

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	vr, _ = vc.BadRequest(fmt.Errorf("err")).(*ViewResult)
	assert.NotNil(vr.ViewModel)
	_, ok := vr.ViewModel.(error)
	assert.True(ok)
}

func TestViewCacheInternalError(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.Nil(vc.Initialize())

	ler, _ := vc.InternalError(ex.Class("only a test")).(*LoggedErrorResult)
	assert.NotNil(ler)
	vr := ler.Result.(*ViewResult)
	assert.Equal(vc.InternalErrorTemplateName, vr.ViewName)
	assert.Equal(http.StatusInternalServerError, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.ViewModel)

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	vr, _ = vc.InternalError(fmt.Errorf("err")).(*ViewResult)
	assert.NotNil(vr.ViewModel)
	_, ok := vr.ViewModel.(error)
	assert.True(ok)
}

func TestViewCacheNotFound(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.Nil(vc.Initialize())

	vr, _ := vc.NotFound().(*ViewResult)
	assert.NotNil(vr)
	assert.Equal(vc.NotFoundTemplateName, vr.ViewName)
	assert.Equal(http.StatusNotFound, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.Nil(vr.ViewModel)

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	vr, _ = vc.NotFound().(*ViewResult)
	assert.NotNil(vr.ViewModel)
	_, ok := vr.ViewModel.(error)
	assert.True(ok)
}

func TestViewCacheNotAuthorized(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.Nil(vc.Initialize())

	vr, _ := vc.NotAuthorized().(*ViewResult)
	assert.NotNil(vr)
	assert.Equal(vc.NotAuthorizedTemplateName, vr.ViewName)
	assert.Equal(http.StatusUnauthorized, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.Nil(vr.ViewModel)

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	vr, _ = vc.NotAuthorized().(*ViewResult)
	assert.NotNil(vr.ViewModel)
	_, ok := vr.ViewModel.(error)
	assert.True(ok)
}

func TestViewCacheStatus(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.Nil(vc.Initialize())

	vr, _ := vc.Status(http.StatusFailedDependency, nil).(*ViewResult)
	assert.NotNil(vr)
	assert.Equal(vc.StatusTemplateName, vr.ViewName)
	assert.Equal(http.StatusFailedDependency, vr.StatusCode)
	assert.NotNil(vr.Views)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.ViewModel)

	vc = MustNewViewCache()
	vc.AddLiterals(`{{define "test"}}failure{{`)
	vr, _ = vc.Status(http.StatusPreconditionFailed, nil).(*ViewResult)
	assert.NotNil(vr.ViewModel)
	_, ok := vr.ViewModel.(error)
	assert.True(ok)
}

func TestViewCacheViewStatus(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	assert.Nil(vc.Initialize())

	vc.AddLiterals(`{{define "test"}}failure{{`)
	vr, _ := vc.ViewStatus(http.StatusPreconditionFailed, "", nil).(*ViewResult)
	assert.NotNil(vr.ViewModel)
	_, ok := vr.ViewModel.(error)
	assert.True(ok)
}

func TestViewCacheView(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()
	vc.AddLiterals(
		`{{ define "test" }}{{.ViewModel}}{{end}}`,
	)
	assert.Nil(vc.Initialize())

	vr, _ := vc.View("test", "foo").(*ViewResult)
	assert.NotNil(vr)
	assert.Equal("test", vr.ViewName)
	assert.Equal(http.StatusOK, vr.StatusCode)
	assert.Equal("foo", vr.ViewModel)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.Views)

	// handle if the view is not found ...
	ler, _ := vc.View("not-test", "foo").(*LoggedErrorResult)
	assert.NotNil(ler)
	vr, _ = ler.Result.(*ViewResult)
	assert.Equal(vc.InternalErrorTemplateName, vr.ViewName)
	assert.Equal(http.StatusInternalServerError, vr.StatusCode)
	assert.NotNil(vr.Template)
	assert.NotNil(vr.Views)
}

func TestViewCacheViewError(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()

	vcr := vc.viewError(ex.Class("test error")).(*ViewResult)
	assert.NotNil(vcr)
	assert.Equal(DefaultTemplateNameInternalError, vcr.ViewName)
	assert.Equal(http.StatusInternalServerError, vcr.StatusCode)
	assert.Equal(ex.Class("test error"), vcr.ViewModel)
	assert.NotNil(vcr.Template)
	assert.NotNil(vcr.Views)
}

func TestViewCacheFuncs(t *testing.T) {
	assert := assert.New(t)

	vc := MustNewViewCache()

	f := func() {}

	vc.FuncMap = nil
	opt := OptViewCacheFunc("useless", f)
	assert.NotNil(opt)
	assert.Nil(opt(vc))
	assert.NotEmpty(vc.FuncMap)
	_, ok := vc.FuncMap["useless"]
	assert.True(ok)

	opt = OptViewCacheFuncMap(template.FuncMap{})
	assert.Nil(opt(vc))
	assert.Empty(vc.FuncMap)
}
