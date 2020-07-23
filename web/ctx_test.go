package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
)

func TestCtxGetState(t *testing.T) {
	assert := assert.New(t)

	context := NewCtx(nil, nil)
	context.WithStateValue("foo", "bar")
	assert.Equal("bar", context.StateValue("foo"))
}

func TestCtxParamQuery(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxQueryValue("foo", "bar"))
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)

	param, err = context.QueryValue("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxParamHeader(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxHeaderValue("foo", "bar"))
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)

	param, err = context.HeaderValue("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxParamForm(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxPostFormValue("foo", "bar"))
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)

	param, err = context.FormValue("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxParamCookie(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxCookieValue("foo", "bar"))
	param, err := context.Param("foo")
	assert.Nil(err)
	assert.Equal("bar", param)
}

func TestCtxPostBodyAsString(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxBodyBytes([]byte("test payload")))
	body, err := context.PostBodyAsString()
	assert.Nil(err)
	assert.Equal("test payload", body)

	context = MockCtx("GET", "/")
	body, err = context.PostBodyAsString()
	assert.Nil(err)
	assert.Empty(body)
}

func TestCtxPostBodyAsJSON(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxBodyBytes([]byte(`{"test":"test payload"}`)))

	var contents map[string]interface{}
	err := context.PostBodyAsJSON(&contents)
	assert.Nil(err)
	assert.Equal("test payload", contents["test"])

	context = MockCtx("GET", "/")
	assert.Nil(err)
	contents = make(map[string]interface{})
	err = context.PostBodyAsJSON(&contents)
	assert.NotNil(err)
}

type postXMLTest string

func TestCtxPostBodyAsXML(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxBodyBytes([]byte(`<postXMLTest>test payload</postXMLTest>`)))

	var contents postXMLTest
	err := context.PostBodyAsXML(&contents)
	assert.Nil(err)
	assert.Equal("test payload", string(contents))
}

func TestCtxPostedFiles(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/")
	postedFiles, err := webutil.PostedFiles(context.Request)
	assert.Nil(err)
	assert.Empty(postedFiles)

	context = MockCtx("GET", "/", OptCtxPostedFiles(webutil.PostedFile{
		Key:      "file",
		FileName: "test.txt",
		Contents: []byte("this is only a test"),
	}))

	postedFiles, err = webutil.PostedFiles(context.Request)
	assert.Nil(err)
	assert.NotEmpty(postedFiles)
	assert.Equal("file", postedFiles[0].Key)
	assert.Equal("test.txt", postedFiles[0].FileName)
	assert.Equal("this is only a test", string(postedFiles[0].Contents))
}

func TestCtxRouteParam(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/", OptCtxRouteParamValue("foo", "bar"))
	value, err := context.RouteParam("foo")
	assert.Nil(err)
	assert.Equal("bar", value)
}

func TestCtxSession(t *testing.T) {
	assert := assert.New(t)

	session := NewSession("test user", NewSessionID())
	ctx := MockCtx("GET", "/", OptCtxSession(session))
	assert.Equal(ctx.Session, session)
}

func TestCtxWriteNewCookie(t *testing.T) {
	assert := assert.New(t)

	context := MockCtx("GET", "/")
	context.WriteNewCookie(&http.Cookie{
		Name:     "foo",
		Value:    "bar",
		Path:     "/foo/bar",
		HttpOnly: true,
		Secure:   true,
	})
	assert.Equal("foo=bar; Path=/foo/bar; Domain=localhost; HttpOnly; Secure", context.Response.Header().Get("Set-Cookie"))
}

func TestCtxExtendCookie(t *testing.T) {
	assert := assert.New(t)

	ctx := MockCtx("GET", "/", OptCtxCookieValue("foo", "bar"))
	ctx.ExtendCookie("foo", "/", 0, 0, 1)

	cookies := ReadSetCookies(ctx.Response.Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.False(cookie.Expires.IsZero())
}

func TestCtxExtendCookieByDuration(t *testing.T) {
	assert := assert.New(t)

	ctx := MockCtx("GET", "/", OptCtxCookieValue("foo", "bar"))
	ctx.ExtendCookieByDuration("foo", "/", time.Hour)

	cookies := ReadSetCookies(ctx.Response.Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.False(cookie.Expires.IsZero())
}

func TestCtxCookieDomain(t *testing.T) {
	assert := assert.New(t)

	// Fallback to `ctx.Request.Host`
	ctx := MockCtx("GET", "/")
	domain := ctx.CookieDomain()
	assert.Equal("localhost", domain)
	assert.Nil(ctx.App)

	// Use `ctx.App.Config.BaseURL`
	cfg := Config{BaseURL: "http://localhost:8080"}
	app := MustNew(OptConfig(cfg))
	ctx = MockCtx("GET", "/", OptCtxApp(app))
	domain = ctx.CookieDomain()
	assert.Equal("localhost", domain)
}

type PostFormTest struct {
	ID       string  `postForm:"id"`
	Name     string  `postForm:"Name"`
	Cost     float64 `postForm:"notCost"`
	Excluded string
}

func TestCtxPostBodyAsForm(t *testing.T) {
	assert := assert.New(t)

	formValues := url.Values{
		"id":       []string{uuid.V4().String()},
		"Name":     []string{"foobar"},
		"notCost":  []string{"3.14", "6.28"},
		"Excluded": []string{"bad"},
	}
	postBody := []byte(formValues.Encode())

	ctx := MockCtx("POST", "/")
	ctx.Request.Header.Set(webutil.HeaderContentType, webutil.ContentTypeApplicationFormEncoded)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(postBody))

	var p PostFormTest
	assert.Nil(ctx.PostBodyAsForm(&p))
	assert.Equal(formValues["id"][0], p.ID)
	assert.Equal(formValues["Name"][0], p.Name)
	assert.Equal(3.14, p.Cost)
	assert.Empty(p.Excluded)
}
