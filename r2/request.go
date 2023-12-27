/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

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
	"github.com/blend/go-sdk/webutil"
)

// New returns a new request.
// The default method is GET.
func New(remoteURL string, options ...Option) *Request {
	var r Request
	u, err := url.Parse(remoteURL)
	if err != nil {
		r.Err = ex.New(err)
		return &r
	}
	u.Host = webutil.RemoveHostEmptyPort(u.Host)
	r.Request = &http.Request{
		Method:		MethodGet,
		URL:		u,
		Proto:		"HTTP/1.1",
		ProtoMajor:	1,
		ProtoMinor:	1,
		Header:		make(http.Header),
		Host:		u.Host,
	}
	for _, option := range options {
		if err = option(&r); err != nil {
			r.Err = err
			return &r
		}
	}
	return &r
}

// Request is a combination of the http.Request options and the underlying client.
type Request struct {
	Request	*http.Request
	// Err is an error set on construction.
	// It is checked before sending the request, and will be returned from any of the
	// methods that execute the request.
	// It is typically set in `New(string,...Option)`.
	Err	error
	// Client is the underlying http client used to make the requests.
	Client	*http.Client
	// Closer is an optional step to run as part of the Close() function.
	Closer	func() error
	// Tracer is used to report span contexts to a distributed tracing collector.
	Tracer	Tracer
	// OnRequest is an array of request lifecycle hooks used for logging
	// or to modify the request on a per call basis before it is sent.
	OnRequest	[]OnRequestListener
	// OnResponse is an array of response lifecycle hooks used typically for logging.
	OnResponse	[]OnResponseListener
}

// WithContext implements the `WithContext` method for the underlying request.
//
// It is preserved here because the pointer indirects are non-trivial.
func (r *Request) WithContext(ctx context.Context) *Request {
	*r.Request = *r.Request.WithContext(ctx)
	return r
}

// Do executes the request.
func (r Request) Do() (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if !webutil.IsValidMethod(r.Request.Method) {
		return nil, ex.New(ErrInvalidMethod, ex.OptMessagef("method: %q", r.Request.Method))
	}

	if len(r.Request.PostForm) > 0 && r.Request.Body == nil {
		body := r.Request.PostForm.Encode()
		buffer := bytes.NewBufferString(body)
		r.Request.ContentLength = int64(buffer.Len())
		r.Request.Body = io.NopCloser(buffer)
	}

	if r.Request.Body == nil {
		r.Request.Body = http.NoBody
		r.Request.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(http.NoBody), nil }
	}

	started := time.Now().UTC()
	var finisher TraceFinisher
	if r.Tracer != nil {
		finisher = r.Tracer.Start(r.Request)
	}
	for _, listener := range r.OnRequest {
		if err := listener(r.Request); err != nil {
			return nil, err
		}
	}

	var err error
	var res *http.Response
	if r.Client != nil {
		res, err = r.Client.Do(r.Request)
	} else {
		res, err = http.DefaultClient.Do(r.Request)
	}
	if finisher != nil {
		finisher.Finish(r.Request, res, started, err)
	}
	for _, listener := range r.OnResponse {
		if listenerErr := listener(r.Request, res, started, err); listenerErr != nil {
			err = ex.Append(err, listenerErr)
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Close closes the request if there is a closer specified.
func (r *Request) Close() error {
	if r.Closer != nil {
		return r.Closer()
	}
	return nil
}

// Discard reads the response fully and discards all data it reads, and returns the response metadata.
func (r Request) Discard() (res *http.Response, err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			err = ex.Append(err, closeErr)
		}
	}()

	res, err = r.Do()
	if err != nil {
		res = nil
		return
	}
	defer res.Body.Close()
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		err = ex.New(err)
		return
	}
	return
}

// CopyTo copies the response body to a given writer.
func (r Request) CopyTo(dst io.Writer) (count int64, err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			err = ex.Append(err, closeErr)
		}
	}()

	var res *http.Response
	res, err = r.Do()
	if err != nil {
		res = nil
		return
	}
	defer res.Body.Close()
	count, err = io.Copy(dst, res.Body)
	if err != nil {
		err = ex.New(err)
		return
	}
	return
}

// Bytes reads the response and returns it as a byte array, along with the response metadata..
func (r Request) Bytes() (contents []byte, res *http.Response, err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			err = ex.Append(err, closeErr)
		}
	}()
	res, err = r.Do()
	if err != nil {
		res = nil
		err = ex.New(err)
		return
	}
	defer func() {
		err = ex.Append(err, res.Body.Close())
	}()
	contents, err = io.ReadAll(res.Body)
	if err != nil {
		err = ex.New(err)
		return
	}
	return
}

// JSON reads the response as json into a given object and returns the response metadata.
func (r Request) JSON(dst interface{}) (res *http.Response, err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			err = ex.Append(err, closeErr)
		}
	}()

	res, err = r.Do()
	if err != nil {
		res = nil
		err = ex.New(err)
		return
	}
	defer func() {
		err = ex.Append(err, res.Body.Close())
	}()
	if res.StatusCode == http.StatusNoContent {
		err = ex.New(ErrNoContentJSON)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(dst); err != nil {
		err = ex.New(err)
		return
	}
	return
}

// JSONBytes reads the response as json into a given object
// and returns the response bytes as well as the response metadata.
//
// This method is useful for debugging responses.
func (r Request) JSONBytes(dst interface{}) (body []byte, res *http.Response, err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			err = ex.Append(err, closeErr)
		}
	}()

	res, err = r.Do()
	if err != nil {
		res = nil
		err = ex.New(err)
		return
	}
	defer func() {
		err = ex.Append(err, res.Body.Close())
	}()
	if res.StatusCode == http.StatusNoContent {
		err = ex.New(ErrNoContentJSON)
		return
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		err = ex.New(err)
		return
	}
	if err = json.Unmarshal(body, dst); err != nil {
		err = ex.New(err)
		return
	}
	return
}

// XML reads the response as xml into a given object and returns the response metadata.
func (r Request) XML(dst interface{}) (res *http.Response, err error) {
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			err = ex.Append(err, closeErr)
		}
	}()

	res, err = r.Do()
	if err != nil {
		res = nil
		err = ex.New(err)
		return
	}
	defer func() {
		err = ex.Append(err, res.Body.Close())
	}()
	if res.StatusCode == http.StatusNoContent {
		err = ex.New(ErrNoContentXML)
		return
	}
	if err = xml.NewDecoder(res.Body).Decode(dst); err != nil {
		err = ex.New(err)
		return
	}
	return
}
