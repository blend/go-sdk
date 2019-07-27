package web

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/webutil"
)

// NewCtx returns a new ctx.
func NewCtx(w ResponseWriter, r *http.Request, options ...CtxOption) *Ctx {
	ctx := Ctx{
		Response: w,
		Request:  r,
		State:    &SyncState{},
	}
	for _, option := range options {
		option(&ctx)
	}
	return &ctx
}

// Ctx is the struct that represents the context for an hc request.
type Ctx struct {
	App             *App
	DefaultProvider ResultProvider
	Views           *ViewCache

	Response    ResponseWriter
	Request     *http.Request
	Tracer      Tracer
	Body        []byte
	Form        url.Values
	State       State
	Session     *Session
	Route       *Route
	RouteParams RouteParameters

	RequestStart time.Time
	RequestEnd   time.Time
}

// WithContext sets the background context for the request.
func (rc *Ctx) WithContext(context context.Context) *Ctx {
	rc.Request = rc.Request.WithContext(context)
	return rc
}

// Context returns the context.
func (rc *Ctx) Context() context.Context {
	return rc.Request.Context()
}

// WithStateValue sets the state for a key to an object.
func (rc *Ctx) WithStateValue(key string, value interface{}) *Ctx {
	rc.State.Set(key, value)
	return rc
}

// StateValue returns an object in the state cache.
func (rc *Ctx) StateValue(key string) interface{} {
	return rc.State.Get(key)
}

// Param returns a parameter from the request.
/*
It checks, in order:
	- RouteParam
	- QueryValue
	- HeaderValue
	- FormValue
	- CookieValue

It should only be used in cases where you don't necessarily know where the param
value will be coming from. Where possible, use the more tightly scoped
param getters.

It returns the value, and a validation error if the value is not found in
any of the possible sources.

You can use one of the Value functions to also cast the resulting string
into a useful type:

	typed, err := web.IntValue(rc.Param("fooID"))

*/
func (rc *Ctx) Param(name string) (value string, err error) {
	if rc.RouteParams != nil {
		value = rc.RouteParams.Get(name)
		if value != "" {
			return
		}
	}
	if rc.Request != nil {
		if rc.Request.URL != nil {
			value = rc.Request.URL.Query().Get(name)
			if value != "" {
				return
			}
		}
		if rc.Request.Header != nil {
			value = rc.Request.Header.Get(name)
			if value != "" {
				return
			}
		}

		value, err = rc.FormValue(name)
		if err == nil {
			return
		}

		var cookie *http.Cookie
		cookie, err = rc.Request.Cookie(name)
		if err == nil && cookie.Value != "" {
			value = cookie.Value
			return
		}
	}

	err = NewParameterMissingError(name)
	return
}

// RouteParam returns a string route parameter
func (rc *Ctx) RouteParam(key string) (output string, err error) {
	if value, hasKey := rc.RouteParams[key]; hasKey {
		output = value
		return
	}
	err = NewParameterMissingError(key)
	return
}

// QueryValue returns a query value.
func (rc *Ctx) QueryValue(key string) (value string, err error) {
	if value = rc.Request.URL.Query().Get(key); len(value) > 0 {
		return
	}
	err = NewParameterMissingError(key)
	return
}

// FormValue returns a form value.
func (rc *Ctx) FormValue(key string) (output string, err error) {
	if err = rc.ensureForm(); err != nil {
		return
	}
	if value := rc.Form.Get(key); len(value) > 0 {
		output = value
		return
	}
	err = NewParameterMissingError(key)
	return
}

// HeaderValue returns a header value.
func (rc *Ctx) HeaderValue(key string) (value string, err error) {
	if value = rc.Request.Header.Get(key); len(value) > 0 {
		return
	}
	err = NewParameterMissingError(key)
	return
}

// PostBody returns the bytes in a post body.
// It will store those bytes for re-use on the context itself.
func (rc *Ctx) PostBody() ([]byte, error) {
	var err error
	if len(rc.Body) == 0 {
		if rc.Request != nil && rc.Request.Body != nil {
			defer rc.Request.Body.Close()
			rc.Body, err = ioutil.ReadAll(rc.Request.Body)
		}
		if err != nil {
			return nil, ex.New(err)
		}
	}
	return rc.Body, nil
}

// PostBodyAsString returns the post body as a string.
func (rc *Ctx) PostBodyAsString() (string, error) {
	body, err := rc.PostBody()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// PostBodyAsJSON reads the incoming post body (closing it) and marshals it to the target object as json.
func (rc *Ctx) PostBodyAsJSON(response interface{}) error {
	body, err := rc.PostBody()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, response); err != nil {
		return ex.New(err)
	}
	return nil
}

// PostBodyAsXML reads the incoming post body (closing it) and marshals it to the target object as xml.
func (rc *Ctx) PostBodyAsXML(response interface{}) error {
	body, err := rc.PostBody()
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(body, response); err != nil {
		return ex.New(err)
	}
	return nil
}

// CookieDomain returns the cookie domain for a request.
func (rc *Ctx) CookieDomain() string {
	if rc.App != nil && rc.App.Config.BaseURL != "" {
		return webutil.MustParseURL(rc.App.Config.BaseURL).Host
	}
	return rc.Request.Host
}

// Cookie returns a named cookie from the request.
func (rc *Ctx) Cookie(name string) *http.Cookie {
	cookie, err := rc.Request.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

// WriteNewCookie is a helper method for WriteCookie.
func (rc *Ctx) WriteNewCookie(cookie *http.Cookie) {
	if cookie.Domain == "" {
		cookie.Domain = rc.CookieDomain()
	}
	http.SetCookie(rc.Response, cookie)
}

// ExtendCookieByDuration extends a cookie by a time duration (on the order of nanoseconds to hours).
func (rc *Ctx) ExtendCookieByDuration(name string, path string, duration time.Duration) {
	c := rc.Cookie(name)
	if c == nil {
		return
	}
	c.Path = path
	c.Domain = rc.CookieDomain()
	if c.Expires.IsZero() {
		c.Expires = time.Now().UTC().Add(duration)
	} else {
		c.Expires = c.Expires.Add(duration)
	}
	http.SetCookie(rc.Response, c)
}

// ExtendCookie extends a cookie by years, months or days.
func (rc *Ctx) ExtendCookie(name string, path string, years, months, days int) {
	c := rc.Cookie(name)
	if c == nil {
		return
	}
	c.Path = path
	c.Domain = rc.CookieDomain()
	if c.Expires.IsZero() {
		c.Expires = time.Now().UTC().AddDate(years, months, days)
	} else {
		c.Expires = c.Expires.AddDate(years, months, days)
	}
	http.SetCookie(rc.Response, c)
}

// ExpireCookie expires a cookie.
func (rc *Ctx) ExpireCookie(name string, path string) {
	c := rc.Cookie(name)
	if c == nil {
		return
	}
	c.Path = path
	c.Value = NewSessionID()
	c.Domain = rc.CookieDomain()
	c.Expires = time.Now().UTC().AddDate(-1, 0, 0)

	http.SetCookie(rc.Response, c)
}

// Elapsed is the time delta between start and end.
func (rc *Ctx) Elapsed() time.Duration {
	if !rc.RequestEnd.IsZero() {
		return rc.RequestEnd.Sub(rc.RequestStart)
	}
	return time.Now().UTC().Sub(rc.RequestStart)
}

// --------------------------------------------------------------------------------
// internal methods
// --------------------------------------------------------------------------------

func (rc *Ctx) ensureForm() error {
	if rc.Form != nil {
		return nil
	}
	if rc.Request.PostForm != nil {
		rc.Form = rc.Request.PostForm
		return nil
	}

	body, err := rc.PostBody()
	if err != nil {
		return err
	}

	r := &http.Request{
		Method: rc.Request.Method,
		Header: rc.Request.Header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	rc.Form = r.PostForm
	return nil
}

func (rc *Ctx) onRequestStart() {
	rc.RequestStart = time.Now().UTC()
}

func (rc *Ctx) onRequestFinish() {
	rc.RequestEnd = time.Now().UTC()
}
