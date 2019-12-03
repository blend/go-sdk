package r2

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/blend/go-sdk/ex"
)

// New returns a new request.
// The default method is GET.
func New(remoteURL string, options ...Option) *Request {
	var r Request
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		r.Err = ex.New(err)
		return &r
	}
	r.Request = http.Request{
		Method: MethodGet,
		URL:    parsedURL,
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
	http.Request

	// Err is an error set on construction.
	// It is checked before sending the request, and will be returned from any of the
	// methods that execute the request.
	// It is typically set in `New(string,...Option)`.
	Err error
	// Client is the underlying http client used to make the requests.
	Client *http.Client
	// Closer is an optional step to run as part of the Close() function.
	Closer func() error
	// Tracer is used to report span contexts to a distributed tracing collector.
	Tracer Tracer
	// OnRequest is an array of request lifecycle hooks used for logging.
	OnRequest []OnRequestListener
	// OnResponse is an array of response lifecycle hooks used for logging.
	OnResponse []OnResponseListener
}

// Do executes the request.
func (r Request) Do() (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	// reconcile post form values
	if r.Request.PostForm != nil && len(r.Request.PostForm) > 0 {
		if r.Request.Body == nil {
			body := r.Request.PostForm.Encode()
			r.Request.Body = ioutil.NopCloser(strings.NewReader(body))
			r.Request.ContentLength = int64(len(body))
		}
	}

	var err error
	started := time.Now().UTC()

	var finisher TraceFinisher
	if r.Tracer != nil {
		finisher = r.Tracer.Start(&r.Request)
	}

	for _, listener := range r.OnRequest {
		if err = listener(&r.Request); err != nil {
			return nil, err
		}
	}

	var res *http.Response
	if r.Client != nil {
		res, err = r.Client.Do(&r.Request)
	} else {
		res, err = http.DefaultClient.Do(&r.Request)
	}
	if finisher != nil {
		finisher.Finish(&r.Request, res, started, err)
	}
	for _, listener := range r.OnResponse {
		if err = listener(&r.Request, res, started, err); err != nil {
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
func (r Request) Discard() (*http.Response, error) {
	defer r.Close()

	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	_, err = io.Copy(ioutil.Discard, res.Body)
	return res, ex.New(err)
}

// CopyTo copies the response body to a given writer.
func (r Request) CopyTo(dst io.Writer) (int64, error) {
	defer r.Close()

	res, err := r.Do()
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	count, err := io.Copy(dst, res.Body)
	if err != nil {
		return count, ex.New(err)
	}
	return count, nil
}

// Bytes reads the response and returns it as a byte array, along with the response metadata..
func (r Request) Bytes() ([]byte, *http.Response, error) {
	defer r.Close()

	res, err := r.Do()
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, ex.New(err)
	}
	return contents, res, nil
}

// JSON reads the response as json into a given object and returns the response metadata.
func (r Request) JSON(dst interface{}) (*http.Response, error) {
	defer r.Close()

	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return res, ex.New(ErrNoContentJSON)
	}
	return res, ex.New(json.NewDecoder(res.Body).Decode(dst))
}

// XML reads the response as xml into a given object and returns the response metadata.
func (r Request) XML(dst interface{}) (*http.Response, error) {
	defer r.Close()

	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return res, ex.New(ErrNoContentXML)
	}
	return res, ex.New(xml.NewDecoder(res.Body).Decode(dst))
}
