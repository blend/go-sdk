package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// Get returns a new get request.
func Get(url string) *Request {
	return New().AsGet().MustWithRawURL(url)
}

// Post returns a new post request with an optional body.
func Post(url string, body []byte) *Request {
	return New().AsPost().MustWithRawURL(url).WithPostBody(body)
}

// New returns a new HTTPRequest instance.
func New() *Request {
	return &Request{
		method:   MethodGet,
		scheme:   "http",
		query:    make(url.Values),
		header:   make(http.Header),
		postData: make(url.Values),
		context:  context.TODO(),
	}
}

// Request makes http requests.
type Request struct {
	log *logger.Logger

	method string
	scheme string
	host   string // host or host:port
	path   string // path (relative paths may omit leading slash)
	query  url.Values

	cookies []*http.Cookie
	header  http.Header

	basicAuthUsername string
	basicAuthPassword string

	contentType string

	postData    url.Values
	postedFiles []PostedFile
	body        []byte

	dialTimeout           time.Duration
	keepAliveTimeout      time.Duration
	responseHeaderTimeout time.Duration
	tlsHandshakeTimeout   time.Duration
	timeout               time.Duration

	tlsClientCert []byte
	tlsClientKey  []byte

	tlsSkipVerify bool
	tlsRootCAPool *x509.CertPool

	keepAlive          *bool
	disableCompression *bool
	transport          *http.Transport

	context context.Context
	trace   *httptrace.ClientTrace
	tracer  Tracer

	state interface{}

	requestStarted  time.Time
	requestHandler  Handler
	responseHandler ResponseHandler
	mockProvider    MockedResponseProvider
}

// WithTracer sets the request tracer.
func (r *Request) WithTracer(tracer Tracer) *Request {
	r.tracer = tracer
	return r
}

// Tracer returns the request tracer.
func (r *Request) Tracer() Tracer {
	return r.tracer
}

// WithRequestHandler configures an event receiver.
func (r *Request) WithRequestHandler(handler Handler) *Request {
	r.requestHandler = handler
	return r
}

// RequestHandler returns the request handler.
func (r *Request) RequestHandler() Handler {
	return r.requestHandler
}

// WithResponseHandler configures an event receiver.
func (r *Request) WithResponseHandler(listener ResponseHandler) *Request {
	r.responseHandler = listener
	return r
}

// ResponseHandler returns the request response handler.
func (r *Request) ResponseHandler() ResponseHandler {
	return r.responseHandler
}

// WithMockProvider mocks a request response.
func (r *Request) WithMockProvider(provider MockedResponseProvider) *Request {
	r.mockProvider = provider
	return r
}

// MockProvider returns the request mock provider.
func (r *Request) MockProvider() MockedResponseProvider {
	return r.mockProvider
}

// WithContext sets a context for the request.
func (r *Request) WithContext(ctx context.Context) *Request {
	r.context = ctx
	return r
}

// Context returns the request's context.
func (r *Request) Context() context.Context {
	return r.context
}

// WithClientTrace sets up a trace for the request.
func (r *Request) WithClientTrace(trace *httptrace.ClientTrace) *Request {
	r.trace = trace
	return r
}

// ClientTrace returns the diagnostics trace object.
func (r *Request) ClientTrace() *httptrace.ClientTrace {
	return r.trace
}

// WithState adds a state object to the request for later usage.
func (r *Request) WithState(state interface{}) *Request {
	r.state = state
	return r
}

// State returns the request state.
func (r *Request) State() interface{} {
	return r.state
}

// WithLogger enables logging with HTTPRequestLogLevelErrors.
func (r *Request) WithLogger(log *logger.Logger) *Request {
	r.log = log
	return r
}

// Logger returns the request diagnostics agent.
func (r *Request) Logger() *logger.Logger {
	return r.log
}

// WithTransport sets a transport for the request.
func (r *Request) WithTransport(transport *http.Transport) *Request {
	r.transport = transport
	return r
}

// Transport returns a shared http transport.
func (r *Request) Transport() *http.Transport {
	return r.transport
}

// WithKeepAlive sets if the request should use the `Connection=keep-alive` header or not.
func (r *Request) WithKeepAlive() *Request {
	r.keepAlive = optBool(true)
	r = r.WithHeader(HeaderConnection, ConnectionKeepAlive)
	return r
}

// KeepAlive returns if the keep alive.
func (r *Request) KeepAlive() bool {
	if r.keepAlive != nil {
		return *r.keepAlive
	}
	return DefaultKeepAlive
}

// WithDisableCompression sets the disable compression value.
func (r *Request) WithDisableCompression(value bool) *Request {
	r.disableCompression = optBool(value)
	return r
}

// DisableCompression returns if the requests transport should disable compression.
func (r *Request) DisableCompression() bool {
	if r.disableCompression != nil {
		return *r.disableCompression
	}
	return false
}

// WithKeepAliveTimeout sets a keep alive timeout for the requests transport.
func (r *Request) WithKeepAliveTimeout(timeout time.Duration) *Request {
	r.keepAliveTimeout = timeout
	return r
}

// KeepAliveTimeout returns the keep alive timeout, ro the time before idle connections are closed.
func (r *Request) KeepAliveTimeout() time.Duration {
	return r.keepAliveTimeout
}

// WithResponseHeaderTimeout sets a timeout
func (r *Request) WithResponseHeaderTimeout(timeout time.Duration) *Request {
	r.responseHeaderTimeout = timeout
	return r
}

// ResponseHeaderTimeout returns a timeout.
func (r *Request) ResponseHeaderTimeout() time.Duration {
	return r.responseHeaderTimeout
}

// WithTLSHandshakeTimeout sets a timeout
func (r *Request) WithTLSHandshakeTimeout(timeout time.Duration) *Request {
	r.tlsHandshakeTimeout = timeout
	return r
}

// TLSHandshakeTimeout returns a timeout.
func (r *Request) TLSHandshakeTimeout() time.Duration {
	return r.tlsHandshakeTimeout
}

// WithContentType sets the `Content-Type` header for the request.
func (r *Request) WithContentType(contentType string) *Request {
	r.contentType = contentType
	return r
}

// ContentType returns the request content type.
func (r *Request) ContentType() string {
	return r.contentType
}

// WithScheme sets the scheme, or protocol, of the request.
func (r *Request) WithScheme(scheme string) *Request {
	r.scheme = scheme
	return r
}

// Scheme returns the request url scheme.
func (r *Request) Scheme() string {
	return r.scheme
}

// WithHost sets the target url host for the request.
func (r *Request) WithHost(host string) *Request {
	r.host = host
	return r
}

// Host returns the host.
func (r *Request) Host() string {
	return r.host
}

// WithPath sets the path component of the host url..
func (r *Request) WithPath(path string) *Request {
	r.path = path
	return r
}

// WithPathf sets the path component of the host url by the format and arguments.
func (r *Request) WithPathf(format string, args ...interface{}) *Request {
	r.path = fmt.Sprintf(format, args...)
	return r
}

// Path returns the request path.
func (r *Request) Path() string {
	return r.path
}

// WithRawURLf sets the url based on a format and args.
func (r *Request) WithRawURLf(format string, args ...interface{}) (*Request, error) {
	return r.WithRawURL(fmt.Sprintf(format, args...))
}

// MustWithRawURLf sets the url based on a format and args.
func (r *Request) MustWithRawURLf(format string, args ...interface{}) *Request {
	return r.MustWithRawURL(fmt.Sprintf(format, args...))
}

// WithRawURL sets the request target url whole hog.
func (r *Request) WithRawURL(rawURL string) (*Request, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return r, err
	}
	return r.WithURL(parsedURL), nil
}

// MustWithRawURL sets the request target url whole hog.
func (r *Request) MustWithRawURL(rawURL string) *Request {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		panic(err)
	}
	return r.WithURL(parsedURL)
}

// WithURL sets the request url target.
func (r *Request) WithURL(target *url.URL) *Request {
	r.scheme = target.Scheme
	r.host = target.Host
	r.path = target.Path
	for key, values := range target.Query() {
		for _, value := range values {
			r.query.Add(key, value)
		}
	}
	return r
}

// URL returns the request target url.
func (r *Request) URL() *url.URL {
	return &url.URL{
		Scheme:   r.scheme,
		Host:     r.host,
		Path:     r.path,
		RawQuery: r.query.Encode(),
	}
}

// WithHeader sets a header on the request.
func (r *Request) WithHeader(field string, value string) *Request {
	r.header.Set(field, value)
	return r
}

// WithHeaders adds a set of headers to the request.
func (r *Request) WithHeaders(headers http.Header) *Request {
	for key, values := range headers {
		for _, value := range values {
			r.header.Set(key, value)
		}
	}
	return r
}

// Header returns the request headers.
func (r *Request) Header() http.Header {
	return r.header
}

// WithQueryString sets a query string value for the host url of the request.
func (r *Request) WithQueryString(field string, value string) *Request {
	r.query.Add(field, value)
	return r
}

// WithCookie sets a cookie for the request.
func (r *Request) WithCookie(cookie *http.Cookie) *Request {
	r.cookies = append(r.cookies, cookie)
	return r
}

// WithPostData sets a post data value for the request.
func (r *Request) WithPostData(field string, value string) *Request {
	r.postData.Add(field, value)
	return r
}

// GetPostData gets the first value set on the request via the WithPostData func associated with the given key.
// If there are no values associated with the key, GetPostData returns the empty string. This is a useful retrieval
// mechanism for mocks
func (r *Request) GetPostData(field string) string {
	return r.postData.Get(field)
}

// WithPostedFile adds a posted file to the multipart form elements of the request.
func (r *Request) WithPostedFile(key, fileName string, fileContents io.Reader) *Request {
	r.postedFiles = append(r.postedFiles, PostedFile{Key: key, FileName: fileName, FileContents: fileContents})
	return r
}

// WithBasicAuth sets the basic auth headers for a request.
func (r *Request) WithBasicAuth(username, password string) *Request {
	r.basicAuthUsername = username
	r.basicAuthPassword = password
	return r
}

// BasicAuth returns the basic auth credentials for the request.
func (r *Request) BasicAuth() (username, password string) {
	username = r.basicAuthUsername
	password = r.basicAuthPassword
	return
}

// WithTimeout sets a timeout for the request.
// This timeout enforces the time between the start of the connection dial to the first response byte.
func (r *Request) WithTimeout(timeout time.Duration) *Request {
	r.timeout = timeout
	return r
}

// Timeout returns the request timeout.
func (r *Request) Timeout() time.Duration {
	return r.timeout
}

// WithDialTimeout sets a dial timeout for the request.
func (r *Request) WithDialTimeout(timeout time.Duration) *Request {
	r.dialTimeout = timeout
	return r
}

// DialTimeout returns the request dial timeout.
func (r *Request) DialTimeout() time.Duration {
	return r.dialTimeout
}

// WithTLSSkipVerify skips the bad certificate checking on TLS requests.
func (r *Request) WithTLSSkipVerify(skipVerify bool) *Request {
	r.tlsSkipVerify = skipVerify
	return r
}

// TLSSkipVerify returns if we should skip server tls verification.
func (r *Request) TLSSkipVerify() bool {
	return r.tlsSkipVerify
}

// WithTLSClientCertPair sets a tls cert on the transport for the request.
func (r *Request) WithTLSClientCertPair(cert, key []byte) *Request {
	r.tlsClientCert = cert
	r.tlsClientKey = key
	return r
}

// WithTLSClientCert sets a tls cert on the transport for the request.
func (r *Request) WithTLSClientCert(cert []byte) *Request {
	r.tlsClientCert = cert
	return r
}

// WithTLSClientKey sets a tls key on the transport for the request.
func (r *Request) WithTLSClientKey(key []byte) *Request {
	r.tlsClientKey = key
	return r
}

// WithTLSRootCAPool sets the root TLS ca pool for the request.
func (r *Request) WithTLSRootCAPool(certPool *x509.CertPool) *Request {
	r.tlsRootCAPool = certPool
	return r
}

// WithMethod sets the http verb/method of the request.
func (r *Request) WithMethod(verb string) *Request {
	r.method = verb
	return r
}

// Method returns the request method.
func (r *Request) Method() string {
	return r.method
}

// AsGet sets the http verb of the request to `GET`.
func (r *Request) AsGet() *Request {
	r.method = MethodGet
	return r
}

// AsPost sets the http verb of the request to `POST`.
func (r *Request) AsPost() *Request {
	r.method = "POST"
	return r
}

// AsPut sets the http verb of the request to `PUT`.
func (r *Request) AsPut() *Request {
	r.method = MethodPut
	return r
}

// AsPatch sets the http verb of the request to `PATCH`.
func (r *Request) AsPatch() *Request {
	r.method = MethodPatch
	return r
}

// AsDelete sets the http verb of the request to `DELETE`.
func (r *Request) AsDelete() *Request {
	r.method = MethodDelete
	return r
}

// AsOptions sets the http verb of the request to `OPTIONS`.
func (r *Request) AsOptions() *Request {
	r.method = MethodOptions
	return r
}

// WithPostBodyAsJSON sets the post body raw to be the json representation of an object.
func (r *Request) WithPostBodyAsJSON(object interface{}) *Request {
	return r.WithPostBodySerialized(object, r.serializeJSON).WithContentType(ContentTypeApplicationJSON)
}

// WithPostBodyAsXML sets the post body raw to be the xml representation of an object.
func (r *Request) WithPostBodyAsXML(object interface{}) *Request {
	return r.WithPostBodySerialized(object, r.serializeXML).WithContentType(ContentTypeApplicationXML)
}

// WithPostBodySerialized sets the post body with the results of the given serializer.
func (r *Request) WithPostBodySerialized(object interface{}, serialize Serializer) *Request {
	body, _ := serialize(object)
	return r.WithPostBody(body)
}

// WithPostBody sets the post body directly.
func (r *Request) WithPostBody(body []byte) *Request {
	r.body = body
	return r
}

// ApplyTransport applies the request settings to a transport.
func (r *Request) ApplyTransport(transport *http.Transport) error {
	if r.responseHeaderTimeout > 0 {
		transport.ResponseHeaderTimeout = r.responseHeaderTimeout
	}
	if r.tlsHandshakeTimeout > 0 {
		transport.TLSHandshakeTimeout = r.tlsHandshakeTimeout
	}
	if r.keepAlive != nil {
		transport.DisableKeepAlives = !*r.keepAlive
	}
	if r.disableCompression != nil {
		transport.DisableCompression = *r.disableCompression
	}
	if r.dialTimeout > 0 || r.keepAliveTimeout > 0 {
		dialer := &net.Dialer{}
		if r.dialTimeout > 0 {
			dialer.Timeout = r.dialTimeout
		}
		if r.keepAliveTimeout > 0 {
			dialer.KeepAlive = r.keepAliveTimeout
		}
		transport.Dial = dialer.Dial
	}
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.RootCAs = r.tlsRootCAPool
	transport.TLSClientConfig.InsecureSkipVerify = r.tlsSkipVerify
	if len(r.tlsClientCert) > 0 && len(r.tlsClientKey) > 0 {
		cert, err := tls.X509KeyPair(r.tlsClientCert, r.tlsClientKey)
		if err != nil {
			return exception.New(err)
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}
	return nil
}

// Meta returns the request as a HTTPRequestMeta.
func (r *Request) Meta() *Meta {
	return &Meta{
		Method:    r.Method(),
		URL:       r.URL(),
		Headers:   r.Headers(),
		StartTime: r.requestStarted,
	}
}

// RequiresTransport returns if there are request settings that require a shared transport.
func (r *Request) RequiresTransport() bool {
	if len(r.tlsClientCert) > 0 && len(r.tlsClientKey) > 0 {
		return true
	}
	if r.tlsSkipVerify {
		return true
	}
	if r.tlsRootCAPool != nil {
		return true
	}
	if r.keepAliveTimeout > 0 {
		return true
	}
	if r.dialTimeout > 0 {
		return true
	}
	if r.tlsHandshakeTimeout > 0 {
		return true
	}

	return false
}

// PostBody returns the current post body.
func (r *Request) PostBody() (io.Reader, error) {
	if len(r.body) > 0 {
		return bytes.NewBuffer(r.body), nil
	} else if len(r.postData) > 0 {
		return bytes.NewBufferString(r.postData.Encode()), nil
	} else if len(r.postedFiles) > 0 {
		body := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(body)
		for _, postedFile := range r.postedFiles {
			partWriter, err := writer.CreateFormFile(postedFile.Key, postedFile.FileName)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(partWriter, postedFile.FileContents)
			if err != nil {
				return nil, err
			}
		}
		r.header.Set(HeaderContentType, writer.FormDataContentType())
		if err := writer.Close(); err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil, nil
}

// Headers returns the headers on the request.
func (r *Request) Headers() http.Header {
	headers := http.Header{}
	for key, values := range r.header {
		for _, value := range values {
			headers.Set(key, value)
		}
	}
	if len(r.contentType) > 0 {
		headers.Set(HeaderContentType, r.contentType)
	} else if len(r.postData) > 0 {
		headers.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
	}
	return headers
}

// Request returns a http.Request for the HTTPRequest.
func (r *Request) Request() (*http.Request, error) {
	body, err := r.PostBody()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.Method(), r.URL().String(), body)
	if err != nil {
		return nil, exception.New(err)
	}
	if r.context != nil {
		req = req.WithContext(r.context)
	}
	if r.trace != nil {
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), r.trace))
	}
	if len(r.basicAuthUsername) > 0 {
		req.SetBasicAuth(r.basicAuthUsername, r.basicAuthPassword)
	}
	if r.cookies != nil {
		for i := 0; i < len(r.cookies); i++ {
			cookie := r.cookies[i]
			req.AddCookie(cookie)
		}
	}
	for key, values := range r.Headers() {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	return req, nil
}

// Response makes the actual request but returns the underlying http.Response object.
func (r *Request) Response() (res *http.Response, err error) {
	var req *http.Request
	req, err = r.Request()
	if err != nil {
		err = exception.New(err)
		return
	}

	if r.tracer != nil {
		tf := r.tracer.Start(req)
		if tf != nil {
			defer func() {
				tf.Finish(req, NewResponseMeta(res), err)
			}()
		}
	}

	r.requestStarted = time.Now().UTC()
	r.logRequest()
	if r.mockProvider != nil {
		mockedRes := r.mockProvider(r)
		if mockedRes != nil {
			res = mockedRes.Response()
			err = mockedRes.Err
			return
		}
	}

	client := &http.Client{}
	if r.RequiresTransport() && r.transport == nil {
		err = exception.New(ErrRequiresTransport)
		return
	}

	if r.transport != nil {
		err = r.ApplyTransport(r.transport)
		if err != nil {
			return
		}
		client.Transport = r.transport
	}
	if r.timeout > 0 {
		client.Timeout = r.timeout
	}

	res, err = client.Do(req)
	if err != nil {
		err = exception.New(err)
	}
	return
}

// Discard executes the request does not pass the response to handlers or events.
func (r *Request) Discard() error {
	_, err := r.DiscardWithMeta()
	return err
}

// DiscardWithMeta discards the response but triggers listeners.
func (r *Request) DiscardWithMeta() (*ResponseMeta, error) {
	res, err := r.Response()
	if err != nil {
		return nil, exception.New(err)
	}
	meta := NewResponseMeta(res)
	if res.Body != nil {
		defer res.Body.Close()
		contentLength, err := io.Copy(ioutil.Discard, res.Body)
		if err != nil {
			return meta, exception.New(err)
		}
		meta.ContentLength = contentLength
		r.logResponse(meta, nil)
	}
	return meta, nil
}

// Execute makes the request and reads the response.
func (r *Request) Execute() error {
	_, err := r.ExecuteWithMeta()
	return err
}

// ExecuteWithMeta makes the request and returns the meta of the response.
func (r *Request) ExecuteWithMeta() (*ResponseMeta, error) {
	res, err := r.Response()
	if err != nil {
		return nil, exception.New(err)
	}
	meta := NewResponseMeta(res)
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, exception.New(err)
		}
		meta.ContentLength = int64(len(contents))
		r.logResponse(meta, contents)
	}

	return meta, nil
}

// Bytes fetches the response as bytes.
func (r *Request) Bytes() ([]byte, error) {
	contents, _, err := r.BytesWithMeta()
	return contents, err
}

// BytesWithMeta fetches the response as bytes with meta.
func (r *Request) BytesWithMeta() ([]byte, *ResponseMeta, error) {
	res, err := r.Response()
	if err != nil {
		return nil, nil, exception.New(err)
	}
	defer res.Body.Close()

	resMeta := NewResponseMeta(res)
	bytes, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, resMeta, exception.New(readErr)
	}

	resMeta.ContentLength = int64(len(bytes))
	r.logResponse(resMeta, bytes)
	return bytes, resMeta, nil
}

// String returns the body of the response as a string.
func (r *Request) String() (string, error) {
	responseStr, _, err := r.StringWithMeta()
	return responseStr, err
}

// StringWithMeta returns the body of the response as a string in addition to the response metadata.
func (r *Request) StringWithMeta() (string, *ResponseMeta, error) {
	contents, meta, err := r.BytesWithMeta()
	return string(contents), meta, err
}

// JSON unmarshals the response as json to an object.
func (r *Request) JSON(destination interface{}) error {
	_, err := r.deserialize(r.jsonDeserializer(destination))
	return err
}

// JSONWithMeta unmarshals the response as json to an object with metadata.
func (r *Request) JSONWithMeta(destination interface{}) (*ResponseMeta, error) {
	return r.deserialize(r.jsonDeserializer(destination))
}

// JSONWithErrorHandler unmarshals the response as json to an object with metadata or an error object depending on the meta.
func (r *Request) JSONWithErrorHandler(successObject interface{}, errorObject interface{}) (*ResponseMeta, error) {
	return r.deserializeWithError(r.jsonDeserializer(successObject), r.jsonDeserializer(errorObject))
}

// JSONError unmarshals the response as json to an object if the meta indiciates an error.
func (r *Request) JSONError(errorObject interface{}) (*ResponseMeta, error) {
	return r.deserializeWithError(nil, r.jsonDeserializer(errorObject))
}

// XML unmarshals the response as xml to an object with metadata.
func (r *Request) XML(destination interface{}) error {
	_, err := r.deserialize(r.xmlDeserializer(destination))
	return err
}

// XMLWithMeta unmarshals the response as xml to an object with metadata.
func (r *Request) XMLWithMeta(destination interface{}) (*ResponseMeta, error) {
	return r.deserialize(r.xmlDeserializer(destination))
}

// XMLWithErrorHandler unmarshals the response as xml to an object with metadata or an error object depending on the meta.
func (r *Request) XMLWithErrorHandler(successObject interface{}, errorObject interface{}) (*ResponseMeta, error) {
	return r.deserializeWithError(r.xmlDeserializer(successObject), r.xmlDeserializer(errorObject))
}

// Deserialized runs a deserializer with the response.
func (r *Request) Deserialized(deserialize Deserializer) (*ResponseMeta, error) {
	meta, responseErr := r.deserialize(func(body []byte) error {
		return deserialize(body)
	})
	return meta, responseErr
}

func (r *Request) deserialize(handler Deserializer) (*ResponseMeta, error) {
	res, err := r.Response()
	if err != nil {
		return nil, exception.New(err)
	}
	defer res.Body.Close()

	meta := NewResponseMeta(res)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return meta, exception.New(err)
	}

	meta.ContentLength = int64(len(body))
	r.logResponse(meta, body)
	if meta.ContentLength > 0 && handler != nil {
		err = handler(body)
	}
	return meta, exception.New(err)
}

func (r *Request) deserializeWithError(okHandler Deserializer, errorHandler Deserializer) (*ResponseMeta, error) {
	res, err := r.Response()
	if err != nil {
		if res != nil {
			// If there is both an error and a response, we make an attempt to build a ResponseMeta
			return safeNewResponseMeta(res), err
		}
		return nil, exception.New(err)
	}

	meta := NewResponseMeta(res)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return meta, exception.New(err)
	}

	meta.ContentLength = int64(len(body))
	r.logResponse(meta, body)
	if meta.ContentLength > 0 {
		if res.StatusCode == http.StatusOK {
			if okHandler != nil {
				err = okHandler(body)
			}
		} else if errorHandler != nil {
			err = errorHandler(body)
		}
	}
	return meta, exception.New(err)
}

func safeNewResponseMeta(res *http.Response) (meta *ResponseMeta) {
	defer func() {
		if r := recover(); r != nil {
			meta = nil
		}
	}()
	meta = NewResponseMeta(res)
	return
}

func (r *Request) logRequest() {
	if r.requestHandler != nil {
		r.requestHandler(r)
	}

	meta := r.Meta()
	if r.log != nil {
		r.log.Trigger(Event{
			ts:  now(),
			req: meta,
		})
	}
}

func (r *Request) logResponse(resMeta *ResponseMeta, responseBody []byte) {
	if r.responseHandler != nil {
		r.responseHandler(r, resMeta, responseBody)
	}

	if r.log != nil {
		reqMeta := r.Meta()
		r.log.Trigger(ResponseEvent{
			ts:   time.Now().UTC(),
			req:  reqMeta,
			res:  resMeta,
			body: responseBody,
		})
	}
}

// Hash returns a hashcode for a request.
func (r *Request) Hash() uint32 {
	buffer := new(bytes.Buffer)
	buffer.WriteString(r.method)
	buffer.WriteRune('|')
	buffer.WriteString(r.URL().String())
	h := fnv.New32a()
	h.Write(buffer.Bytes())
	return h.Sum32()
}

// Equals returns if a request equals another request.
func (r *Request) Equals(other *Request) bool {
	if other == nil {
		return false
	}
	if r.method != other.method {
		return false
	}
	if r.URL().String() != other.URL().String() {
		return false
	}
	return true
}

//--------------------------------------------------------------------------------
// Unexported Utility Functions
//--------------------------------------------------------------------------------

func (r *Request) jsonDeserializer(object interface{}) Deserializer {
	return func(body []byte) error {
		return r.deserializeJSON(object, body)
	}
}

func (r *Request) xmlDeserializer(object interface{}) Deserializer {
	return func(body []byte) error {
		return r.deserializeXML(object, body)
	}
}

func (r *Request) deserializeJSON(object interface{}, body []byte) error {
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decodeErr := decoder.Decode(object)
	return exception.New(decodeErr)
}

func (r *Request) deserializeJSONFromReader(object interface{}, body io.Reader) error {
	decoder := json.NewDecoder(body)
	decodeErr := decoder.Decode(object)
	return exception.New(decodeErr)
}

func (r *Request) serializeJSON(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

func (r *Request) serializeJSONToReader(object interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(object)
	return buf, err
}

func (r *Request) deserializeXML(object interface{}, body []byte) error {
	return r.deserializeXMLFromReader(object, bytes.NewBuffer(body))
}

func (r *Request) deserializeXMLFromReader(object interface{}, reader io.Reader) error {
	decoder := xml.NewDecoder(reader)
	return decoder.Decode(object)
}

func (r *Request) serializeXML(object interface{}) ([]byte, error) {
	return xml.Marshal(object)
}

func (r *Request) serializeXMLToReader(object interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := xml.NewEncoder(buf)
	err := encoder.Encode(object)
	return buf, err
}
