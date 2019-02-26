package web

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/webutil"
)

// NewCtxID returns a pseudo-unique key for a context.
func NewCtxID() string {
	return stringutil.Random(stringutil.Letters, 12)
}

// NewCtx returns a new ctx.
func NewCtx(w ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		id:       NewCtxID(),
		response: w,
		request:  r,
		state:    &SyncState{},
	}
}

// NewMockCtx returns a new mock ctx.
// It is intended to be used in testing.
func NewMockCtx(method, path string) *Ctx {
	return NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest(method, path))
}

// Ctx is the struct that represents the context for an hc request.
type Ctx struct {
	id       string
	response ResponseWriter
	request  *http.Request

	app   *App
	views *ViewCache
	log   logger.FullReceiver
	auth  *AuthManager

	tracer Tracer

	postBody []byte
	form     url.Values

	defaultResultProvider ResultProvider

	state       State
	session     *Session
	route       *Route
	routeParams RouteParameters

	requestStart time.Time
	requestEnd   time.Time
}

// WithID sets the context ID.
func (rc *Ctx) WithID(id string) *Ctx {
	rc.id = id
	return rc
}

// ID returns a pseudo unique id for the request.
func (rc *Ctx) ID() string {
	return rc.id
}

// WithResponse sets the underlying response.
func (rc *Ctx) WithResponse(res ResponseWriter) *Ctx {
	rc.response = res
	return rc
}

// Response returns the underyling response.
func (rc *Ctx) Response() ResponseWriter {
	return rc.response
}

// WithRequest sets the underlying request.
func (rc *Ctx) WithRequest(req *http.Request) *Ctx {
	rc.request = req
	return rc
}

// Request returns the underlying request.
func (rc *Ctx) Request() *http.Request {
	return rc.request
}

// WithTracer sets the tracer.
func (rc *Ctx) WithTracer(tracer Tracer) *Ctx {
	rc.tracer = tracer
	return rc
}

// Tracer returns the tracer.
func (rc *Ctx) Tracer() Tracer {
	return rc.tracer
}

// WithContext sets the background context for the request.
func (rc *Ctx) WithContext(context context.Context) *Ctx {
	rc.request = rc.request.WithContext(context)
	return rc
}

// Context returns the context.
func (rc *Ctx) Context() context.Context {
	return rc.request.Context()
}

// WithApp sets the app reference for the ctx.
func (rc *Ctx) WithApp(app *App) *Ctx {
	rc.app = app
	return rc
}

// App returns the app reference.
func (rc *Ctx) App() *App {
	return rc.app
}

// WithAuth sets the request context auth.
func (rc *Ctx) WithAuth(authManager *AuthManager) *Ctx {
	rc.auth = authManager
	return rc
}

// Auth returns the AuthManager for the request.
func (rc *Ctx) Auth() *AuthManager {
	return rc.auth
}

// WithSession sets the session for the request.
func (rc *Ctx) WithSession(session *Session) *Ctx {
	rc.session = session
	return rc
}

// Session returns the session (if any) on the request.
func (rc *Ctx) Session() *Session {
	return rc.session
}

// WithViews sets the view cache for the ctx.
func (rc *Ctx) WithViews(vc *ViewCache) *Ctx {
	rc.views = vc
	return rc
}

// View returns the view cache as a result provider.
/*
It returns a reference to the view cache where views can either be read from disk
for every request (uncached) or read from an in-memory cache.

To return a web result for a view with the name "index" simply return:

	return r.View().View("index", myViewmodel)

It is important to not you'll want to have loaded the "index" view at some point
in the application bootstrap (typically when you register your controller).
*/
func (rc *Ctx) View() *ViewCache {
	return rc.views
}

// JSON returns the JSON result provider.
/*
It can be eschewed for:

	return web.JSON.Result(foo)

But is left in place for legacy reasons.
*/
func (rc *Ctx) JSON() JSONResultProvider {
	return JSON
}

// XML returns the xml result provider.
/*
It can be eschewed for:

	return web.XML.Result(foo)

But is left in place for legacy reasons.
*/
func (rc *Ctx) XML() XMLResultProvider {
	return XML
}

// Text returns the text result provider.
/*
It can be eschewed for:

	return web.Text.Result(foo)

But is left in place for legacy reasons.
*/
func (rc *Ctx) Text() TextResultProvider {
	return Text
}

// DefaultResultProvider returns the current result provider for the context. This is
// set by calling SetDefaultResultProvider or using one of the pre-built middleware
// steps that set it for you.
func (rc *Ctx) DefaultResultProvider() ResultProvider {
	return rc.defaultResultProvider
}

// WithDefaultResultProvider sets the default result provider.
func (rc *Ctx) WithDefaultResultProvider(provider ResultProvider) *Ctx {
	rc.defaultResultProvider = provider
	return rc
}

// WithState sets the state.
func (rc *Ctx) WithState(state State) *Ctx {
	rc.state = state
	return rc
}

// State returns the state.
func (rc *Ctx) State() State {
	return rc.state
}

// WithStateValue sets the state for a key to an object.
func (rc *Ctx) WithStateValue(key string, value interface{}) *Ctx {
	rc.state.Set(key, value)
	return rc
}

// StateValue returns an object in the state cache.
func (rc *Ctx) StateValue(key string) interface{} {
	return rc.state.Get(key)
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
func (rc *Ctx) Param(name string) (string, error) {
	if rc.routeParams != nil {
		routeValue := rc.routeParams.Get(name)
		if len(routeValue) > 0 {
			return routeValue, nil
		}
	}
	if rc.request != nil {
		if rc.request.URL != nil {
			queryValue := rc.request.URL.Query().Get(name)
			if len(queryValue) > 0 {
				return queryValue, nil
			}
		}
		if rc.request.Header != nil {
			headerValue := rc.request.Header.Get(name)
			if len(headerValue) > 0 {
				return headerValue, nil
			}
		}

		formValue, err := rc.FormValue(name)
		if err == nil {
			return formValue, nil
		}

		cookie, cookieErr := rc.request.Cookie(name)
		if cookieErr == nil && len(cookie.Value) != 0 {
			return cookie.Value, nil
		}
	}

	return "", NewParameterMissingError(name)
}

// RouteParam returns a string route parameter
func (rc *Ctx) RouteParam(key string) (output string, err error) {
	if value, hasKey := rc.routeParams[key]; hasKey {
		output = value
		return
	}
	err = NewParameterMissingError(key)
	return
}

// QueryValue returns a query value.
func (rc *Ctx) QueryValue(key string) (value string, err error) {
	if value = rc.request.URL.Query().Get(key); len(value) > 0 {
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
	if value := rc.form.Get(key); len(value) > 0 {
		output = value
		return
	}
	err = NewParameterMissingError(key)
	return
}

// HeaderValue returns a header value.
func (rc *Ctx) HeaderValue(key string) (value string, err error) {
	if value = rc.request.Header.Get(key); len(value) > 0 {
		return
	}
	err = NewParameterMissingError(key)
	return
}

// PostBody returns the bytes in a post body.
func (rc *Ctx) PostBody() ([]byte, error) {
	var err error
	if len(rc.postBody) == 0 {
		if rc.request != nil && rc.request.Body != nil {
			defer rc.request.Body.Close()
			rc.postBody, err = ioutil.ReadAll(rc.request.Body)
		}
		if err != nil {
			return nil, exception.New(err)
		}
	}
	return rc.postBody, nil
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
		return exception.New(err)
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
		return exception.New(err)
	}
	return nil
}

// PostedFiles returns any files posted
func (rc *Ctx) PostedFiles() ([]PostedFile, error) {
	var files []PostedFile

	err := rc.request.ParseMultipartForm(PostBodySize)
	if err == nil {
		for key := range rc.request.MultipartForm.File {
			fileReader, fileHeader, err := rc.request.FormFile(key)
			if err != nil {
				return nil, exception.New(err)
			}
			bytes, err := ioutil.ReadAll(fileReader)
			if err != nil {
				return nil, exception.New(err)
			}
			files = append(files, PostedFile{Key: key, FileName: fileHeader.Filename, Contents: bytes})
		}
	} else {
		err = rc.request.ParseForm()
		if err == nil {
			for key := range rc.request.PostForm {
				if fileReader, fileHeader, err := rc.request.FormFile(key); err == nil && fileReader != nil {
					bytes, err := ioutil.ReadAll(fileReader)
					if err != nil {
						return nil, exception.New(err)
					}
					files = append(files, PostedFile{Key: key, FileName: fileHeader.Filename, Contents: bytes})
				}
			}
		}
	}
	return files, nil
}

// GetCookie returns a named cookie from the request.
func (rc *Ctx) GetCookie(name string) *http.Cookie {
	cookie, err := rc.request.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

// WriteCookie writes the cookie to the response.
func (rc *Ctx) WriteCookie(cookie *http.Cookie) {
	http.SetCookie(rc.response, cookie)
}

// WriteNewCookie is a helper method for WriteCookie.
func (rc *Ctx) WriteNewCookie(name string, value string, expires time.Time, path string, secure bool) {
	c := &http.Cookie{
		Name:     name,
		HttpOnly: true, // this is always on because javascript is bad.
		Value:    value,
		Path:     path,
		Secure:   secure,
		Domain:   rc.getCookieDomain(),
		Expires:  expires,
	}
	rc.WriteCookie(c)

}

// ExtendCookieByDuration extends a cookie by a time duration (on the order of nanoseconds to hours).
func (rc *Ctx) ExtendCookieByDuration(name string, path string, duration time.Duration) {
	c := rc.GetCookie(name)
	if c == nil {
		return
	}
	c.Path = path
	c.Domain = rc.getCookieDomain()
	if c.Expires.IsZero() {
		c.Expires = time.Now().UTC().Add(duration)
	} else {
		c.Expires = c.Expires.Add(duration)
	}
	rc.WriteCookie(c)
}

// ExtendCookie extends a cookie by years, months or days.
func (rc *Ctx) ExtendCookie(name string, path string, years, months, days int) {
	c := rc.GetCookie(name)
	if c == nil {
		return
	}
	c.Path = path
	c.Domain = rc.getCookieDomain()
	if c.Expires.IsZero() {
		c.Expires = time.Now().UTC().AddDate(years, months, days)
	} else {
		c.Expires = c.Expires.AddDate(years, months, days)
	}
	rc.WriteCookie(c)
}

// ExpireCookie expires a cookie.
func (rc *Ctx) ExpireCookie(name string, path string) {
	c := rc.GetCookie(name)
	if c == nil {
		return
	}
	c.Path = path
	c.Value = NewSessionID()
	c.Domain = rc.getCookieDomain()
	c.Expires = time.Now().UTC().AddDate(-1, 0, 0)
	rc.WriteCookie(c)
}

// --------------------------------------------------------------------------------
// Logger
// --------------------------------------------------------------------------------

// WithLogger sets the logger.
func (rc *Ctx) WithLogger(log logger.FullReceiver) *Ctx {
	rc.log = log
	return rc
}

// Logger returns the diagnostics agent.
func (rc *Ctx) Logger() logger.FullReceiver {
	return rc.log
}

// --------------------------------------------------------------------------------
// Basic result providers
// --------------------------------------------------------------------------------

// Raw returns a binary response body, sniffing the content type.
func (rc *Ctx) Raw(body []byte) *RawResult {
	return rc.RawWithContentType(http.DetectContentType(body), body)
}

// RawWithContentType returns a binary response with a given content type.
func (rc *Ctx) RawWithContentType(contentType string, body []byte) *RawResult {
	return &RawResult{ContentType: contentType, Response: body}
}

// NoContent returns a service response.
func (rc *Ctx) NoContent() NoContentResult {
	return NoContent
}

// Static returns a static result.
func (rc *Ctx) Static(filePath string) *StaticResult {
	return NewStaticResultForFile(filePath)
}

// Redirect returns a redirect result to a given destination.
func (rc *Ctx) Redirect(destination string) *RedirectResult {
	return &RedirectResult{
		RedirectURI: destination,
	}
}

// Redirectf returns a redirect result to a given destination specified by a given format and scan arguments.
func (rc *Ctx) Redirectf(format string, args ...interface{}) *RedirectResult {
	return &RedirectResult{
		RedirectURI: fmt.Sprintf(format, args...),
	}
}

// RedirectWithMethod returns a redirect result to a destination with a given method.
func (rc *Ctx) RedirectWithMethod(method, destination string) *RedirectResult {
	return &RedirectResult{
		Method:      method,
		RedirectURI: destination,
	}
}

// RedirectWithMethodf returns a redirect result to a destination composed of a format and scan arguments with a given method.
func (rc *Ctx) RedirectWithMethodf(method, format string, args ...interface{}) *RedirectResult {
	return &RedirectResult{
		Method:      method,
		RedirectURI: fmt.Sprintf(format, args...),
	}
}

// Start returns the start request time.
func (rc Ctx) Start() time.Time {
	return rc.requestStart
}

// End returns the end request time.
func (rc Ctx) End() time.Time {
	return rc.requestEnd
}

// Elapsed is the time delta between start and end.
func (rc *Ctx) Elapsed() time.Duration {
	if !rc.requestEnd.IsZero() {
		return rc.requestEnd.Sub(rc.requestStart)
	}
	return time.Now().UTC().Sub(rc.requestStart)
}

// WithRoute sets the route.
func (rc *Ctx) WithRoute(route *Route) *Ctx {
	rc.route = route
	return rc
}

// Route returns the original route match for the request.
func (rc *Ctx) Route() *Route {
	return rc.route
}

// WithRouteParams sets the route parameters.
func (rc *Ctx) WithRouteParams(params RouteParameters) *Ctx {
	rc.routeParams = params
	return rc
}

// RouteParams returns the route parameters for the request.
func (rc *Ctx) RouteParams() RouteParameters {
	return rc.routeParams
}

// --------------------------------------------------------------------------------
// internal methods
// --------------------------------------------------------------------------------

func (rc *Ctx) ensureForm() error {
	if rc.form != nil {
		return nil
	}
	if rc.Request().Form != nil {
		rc.form = rc.Request().Form
		return nil
	}

	body, err := rc.PostBody()
	if err != nil {
		return err
	}

	r := &http.Request{
		Method: rc.Request().Method,
		Header: rc.Request().Header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	rc.form = r.Form
	return nil
}

func (rc *Ctx) getCookieDomain() string {
	if rc.app != nil && rc.app.baseURL != nil {
		return rc.app.baseURL.Host
	}
	return rc.request.Host
}

func (rc *Ctx) onRequestStart() {
	rc.requestStart = time.Now().UTC()
}

func (rc *Ctx) onRequestFinish() {
	rc.requestEnd = time.Now().UTC()
}
