package webutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/ex"
)

var (
	// ErrURLUnset is a (hopefully) uncommon error.
	ErrURLUnset = ex.Class("request url unset")

	// DefaultRequestTimeout is the default webhook timeout.
	DefaultRequestTimeout = 10 * time.Second

	// DefaultRequestMethod is the default webhook method.
	DefaultRequestMethod = "POST"
)

// NewRequestSender creates a new request sender.
/*
A request sender is a sepcialized request factory that makes request to a single endpoint.

It is useful when calling out to predefined things like webhooks.

You can send either raw bytes as the contents.
*/
func NewRequestSender(destination *url.URL) *RequestSender {
	transport := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
	}
	return &RequestSender{
		URL:       destination,
		method:    DefaultRequestMethod,
		transport: transport,
		headers:   http.Header{},
		client: &http.Client{
			Transport: transport,
			Timeout:   DefaultRequestTimeout,
		},
	}
}

// RequestSender is a slack webhook sender.
type RequestSender struct {
	*url.URL

	transport *http.Transport
	close     bool
	method    string
	client    *http.Client
	headers   http.Header
	tracer    RequestTracer
}

// Send sends a request to the destination without a payload.
func (rs *RequestSender) Send() (*http.Response, error) {
	return rs.send(rs.req())
}

// SendBytes sends a message to the webhook with a given msg body as raw bytes.
func (rs *RequestSender) SendBytes(ctx context.Context, contents []byte) (*http.Response, error) {
	req, err := rs.reqBytes(contents)
	if err != nil {
		return nil, err
	}
	return rs.send(req.WithContext(ctx))
}

// SendJSON sends a message to the webhook with a given msg body as json.
func (rs *RequestSender) SendJSON(ctx context.Context, contents interface{}) (*http.Response, error) {
	req, err := rs.reqJSON(contents)
	if err != nil {
		return nil, err
	}
	return rs.send(req.WithContext(ctx))
}

// properties

// WithMethod sets the request method (defaults to POST).
func (rs *RequestSender) WithMethod(method string) *RequestSender {
	rs.method = method
	return rs
}

// Method is the request method.
// It defaults to "POST".
func (rs *RequestSender) Method() string {
	if len(rs.method) == 0 {
		return DefaultRequestMethod
	}
	return rs.method
}

// WithTracer sets the request tracer.
func (rs *RequestSender) WithTracer(tracer RequestTracer) *RequestSender {
	rs.tracer = tracer
	return rs
}

// Tracer returns the request tracer.
func (rs *RequestSender) Tracer() RequestTracer {
	return rs.tracer
}

// WithClose sets if we should close the connection.
func (rs *RequestSender) WithClose(close bool) *RequestSender {
	rs.close = close
	rs.transport.DisableKeepAlives = close
	return rs
}

// Close returns if we should close the connection.
func (rs *RequestSender) Close() bool {
	return rs.close
}

// WithHeaders sets headers.
func (rs *RequestSender) WithHeaders(headers http.Header) *RequestSender {
	rs.headers = headers
	return rs
}

// WithHeader adds an individual header.
func (rs *RequestSender) WithHeader(key, value string) *RequestSender {
	rs.headers.Set(key, value)
	return rs
}

// Headers returns the headers.
func (rs *RequestSender) Headers() http.Header {
	return rs.headers
}

// WithTransport sets the transport.
func (rs *RequestSender) WithTransport(transport *http.Transport) *RequestSender {
	rs.transport = transport
	return rs
}

// Transport returns the transport.
func (rs *RequestSender) Transport() *http.Transport {
	return rs.transport
}

// Client returns the underlying client.
func (rs *RequestSender) Client() *http.Client {
	return rs.client
}

// internal methods

// Send sends a message to the webhook with a given msg body as raw bytes.
func (rs *RequestSender) send(req *http.Request) (res *http.Response, err error) {
	if req.URL == nil {
		return nil, ErrURLUnset
	}
	if rs.tracer != nil {
		tf := rs.tracer.Start(req)
		if tf != nil {
			defer func() { tf.Finish(req, res, err) }()
		}
	}
	res, err = rs.client.Do(req)
	return
}

func (rs *RequestSender) req() *http.Request {
	return &http.Request{
		Method: rs.Method(),
		Close:  rs.Close(),
		URL:    rs.URL,
		Header: rs.headers,
	}
}

func (rs *RequestSender) reqBytes(contents []byte) (*http.Request, error) {
	req := rs.req()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
	req.ContentLength = int64(len(contents))
	return req, nil
}

func (rs *RequestSender) reqJSON(msg interface{}) (*http.Request, error) {
	contents, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req := rs.req()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(contents))
	req.ContentLength = int64(len(contents))
	req.Header.Set(HeaderContentType, ContentTypeApplicationJSON)
	return req, nil
}
