/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/reflectutil"
)

var (
	_ io.Closer = (*Ctx)(nil)
)

// NewCtx returns a new ctx.
func NewCtx(w ResponseWriter, r *http.Request, options ...CtxOption) *Ctx {
	ctx := Ctx{
		Response:	w,
		Request:	r,
		State:		new(SyncState),
	}
	for _, option := range options {
		option(&ctx)
	}
	return &ctx
}

// Ctx is the struct that represents the context for an hc request.
type Ctx struct {
	// App is a reference back to the parent application.
	App	*App
	// Auth is a reference to the app default auth manager, but
	// can be overwritten by middleware.
	Auth	AuthManager
	// DefaultProvider is the app default result provider by default
	// but can be overwritten by middleware.
	DefaultProvider	ResultProvider
	// Views is the app view cache by default but can be
	// overwritten by middleware.
	Views	*ViewCache
	// Response is the response writer for the request.
	Response	ResponseWriter
	// Request is the inbound request metadata.
	Request	*http.Request
	// Body is a cached copy of the post body of a request.
	// It is typically set by calling `.PostBody()` on this context.
	// If you're expecting a large post body, do not use
	// the `.PostBody()` function, instead read directly from `.Request.Body` with
	// a stream reader or similar.
	Body	[]byte
	// Form is a cache of parsed url form values from the post body.
	Form	url.Values
	// State is a mutable bag of state, it contains by default
	// state set on the application.
	State	State
	// Session is the current auth session
	Session	*Session
	// Route is the matching route for the request if relevant.
	Route	*Route
	// RouteParams is a cache of parameters or variables
	// within the route and their values.
	RouteParams	RouteParameters
	// Log is the request specific logger.
	Log	logger.Log
	// Tracer is the app tracer by default if one is set.
	// It can be overwritten by middleware.
	Tracer	Tracer
	// RequestStarted is the time the request was received.
	RequestStarted	time.Time
}

// Close closes the context.
func (rc *Ctx) Close() error {
	if rc.Response != nil {
		if err := rc.Response.Close(); err != nil {
			return err
		}
	}
	return nil
}

// WithContext sets the background context for the request.
func (rc *Ctx) WithContext(ctx context.Context) *Ctx {
	*rc.Request = *rc.Request.WithContext(ctx)
	return rc
}

// Context returns the context.
func (rc *Ctx) Context() context.Context {
	ctx := logger.WithLabels(rc.Request.Context(), logger.GetLabels(rc.Request.Context()))
	ctx = logger.WithLabels(ctx, rc.Labels())
	ctx = logger.WithAnnotations(ctx, logger.CombineAnnotations(logger.GetAnnotations(rc.Request.Context()), rc.Annotations()))
	return ctx
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
	if err = rc.EnsureForm(); err != nil {
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

// PostBody reads, caches and returns the bytes on a request post body.
// It will store those bytes for re-use on this context object.
// If you're expecting a large post body, or a large post body is even possible
// use a stream reader on `.Request.Body` instead of this method.
func (rc *Ctx) PostBody() ([]byte, error) {
	if len(rc.Body) == 0 {
		if rc.Request != nil && rc.Request.GetBody != nil {
			reader, err := rc.Request.GetBody()
			if err != nil {
				return nil, ex.New(err)
			}
			defer reader.Close()
			rc.Body, err = io.ReadAll(reader)
			if err != nil {
				return nil, ex.New(err)
			}
		}
		if rc.Request != nil && rc.Request.Body != nil {
			defer rc.Request.Body.Close()
			var err error
			rc.Body, err = io.ReadAll(rc.Request.Body)
			if err != nil {
				return nil, ex.New(err)
			}
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

// PostBodyAsForm reads the incoming post body (closing it) sets a given object from the post form fields.
// NOTE: the request method *MUST* not be `GET` otherwise the golang internals will skip parsing the body.
func (rc *Ctx) PostBodyAsForm(response interface{}) error {
	if err := rc.EnsureForm(); err != nil {
		return err
	}
	return reflectutil.PatchStringsFunc("postForm", func(key string) (string, bool) {
		if values, ok := rc.Form[key]; ok {
			if len(values) > 0 {
				return values[0], true
			}
			return "", false
		}
		return "", false
	}, response)
}

// Cookie returns a named cookie from the request.
func (rc *Ctx) Cookie(name string) *http.Cookie {
	cookie, err := rc.Request.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

// ExtendCookieByDuration extends a cookie by a time duration (on the order of nanoseconds to hours).
func (rc *Ctx) ExtendCookieByDuration(name string, path string, duration time.Duration) {
	c := rc.Cookie(name)
	if c == nil {
		return
	}
	c.Path = path
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
	// c.MaxAge<0 means delete cookie now, and is equivalent to
	// the literal cookie header content 'Max-Age: 0'
	c.MaxAge = -1
	http.SetCookie(rc.Response, c)
}

// Elapsed is the time delta between start and end.
func (rc *Ctx) Elapsed() time.Duration {
	return time.Now().UTC().Sub(rc.RequestStarted)
}

// --------------------------------------------------------------------------------
// internal methods
// --------------------------------------------------------------------------------

// EnsureForm parses the post body as an application form.
// The parsed form will be available on the `.Form` field.
func (rc *Ctx) EnsureForm() error {
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
		Method:	rc.Request.Method,
		Header:	rc.Request.Header,
		Body:	io.NopCloser(bytes.NewBuffer(body)),
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	rc.Form = r.PostForm
	return nil
}

// Labels returns the labels for logging calls.
func (rc *Ctx) Labels() map[string]string {
	fields := make(map[string]string)
	if rc.Route != nil {
		fields["web.route"] = rc.Route.String()
	}
	if rc.Session != nil {
		fields["web.user"] = rc.Session.UserID
	}
	return fields
}

// Annotations returns the annotations for logging calls.
func (rc *Ctx) Annotations() map[string]interface{} {
	fields := make(map[string]interface{})
	if len(rc.RouteParams) > 0 {
		fields["web.route_parameters"] = rc.RouteParams
	}
	return fields
}
