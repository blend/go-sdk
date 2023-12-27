/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

// RequestOption is an option for http.Request.
type RequestOption func(*http.Request) error

// RequestOptions are an array of RequestOption.
type RequestOptions []RequestOption

// Apply applies the options to a request.
func (ro RequestOptions) Apply(req *http.Request) (err error) {
	for _, option := range ro {
		if err = option(req); err != nil {
			return
		}
	}
	return
}

// OptMethod sets the request method.
func OptMethod(method string) RequestOption {
	return func(r *http.Request) error {
		r.Method = method
		return nil
	}
}

// OptGet sets the request method.
func OptGet() RequestOption {
	return func(r *http.Request) error {
		r.Method = "GET"
		return nil
	}
}

// OptPost sets the request method.
func OptPost() RequestOption {
	return func(r *http.Request) error {
		r.Method = "POST"
		return nil
	}
}

// OptPut sets the request method.
func OptPut() RequestOption {
	return func(r *http.Request) error {
		r.Method = "PUT"
		return nil
	}
}

// OptPatch sets the request method.
func OptPatch() RequestOption {
	return func(r *http.Request) error {
		r.Method = "PATCH"
		return nil
	}
}

// OptDelete sets the request method.
func OptDelete() RequestOption {
	return func(r *http.Request) error {
		r.Method = "DELETE"
		return nil
	}
}

// OptContext sets the request context.
func OptContext(ctx context.Context) RequestOption {
	return func(r *http.Request) error {
		*r = *r.WithContext(ctx)
		return nil
	}
}

// OptBasicAuth is an option that sets the http basic auth.
func OptBasicAuth(username, password string) RequestOption {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.SetBasicAuth(username, password)
		return nil
	}
}

// OptQuery sets the full querystring.
func OptQuery(query url.Values) RequestOption {
	return func(r *http.Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.RawQuery = query.Encode()
		return nil
	}
}

// OptQueryValue sets a query value on a request.
func OptQueryValue(key, value string) RequestOption {
	return func(r *http.Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		existing := r.URL.Query()
		existing.Set(key, value)
		r.URL.RawQuery = existing.Encode()
		return nil
	}
}

// OptQueryValueAdd adds a query value on a request.
func OptQueryValueAdd(key, value string) RequestOption {
	return func(r *http.Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		existing := r.URL.Query()
		existing.Add(key, value)
		r.URL.RawQuery = existing.Encode()
		return nil
	}
}

// OptHeader sets the request headers.
func OptHeader(headers http.Header) RequestOption {
	return func(r *http.Request) error {
		r.Header = headers
		return nil
	}
}

// OptHeaderValue sets a header value on a request.
func OptHeaderValue(key, value string) RequestOption {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set(key, value)
		return nil
	}
}

// OptPostForm sets the request post form and the content type.
func OptPostForm(postForm url.Values) RequestOption {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
		r.PostForm = postForm
		return nil
	}
}

// OptPostFormValue sets a request post form value.
func OptPostFormValue(key, value string) RequestOption {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationFormEncoded)
		if r.PostForm == nil {
			r.PostForm = url.Values{}
		}
		r.PostForm.Set(key, value)
		return nil
	}
}

// OptCookie adds a cookie to a context.
func OptCookie(cookie *http.Cookie) RequestOption {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.AddCookie(cookie)
		return nil
	}
}

// OptCookieValue adds a cookie value to a context.
func OptCookieValue(key, value string) RequestOption {
	return OptCookie(&http.Cookie{Name: key, Value: value})
}

// OptBody sets the post body on the request.
func OptBody(contents io.ReadCloser) RequestOption {
	return func(r *http.Request) error {
		r.Body = contents
		return nil
	}
}

// OptBodyBytes sets a body on a context from bytes.
func OptBodyBytes(body []byte) RequestOption {
	return func(r *http.Request) error {
		r.ContentLength = int64(len(body))
		r.Body = io.NopCloser(bytes.NewReader(body))
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
		r.ContentLength = int64(len(body))
		return nil
	}
}

// OptPostedFiles sets a body from posted files.
func OptPostedFiles(files ...PostedFile) RequestOption {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}

		b := new(bytes.Buffer)
		w := multipart.NewWriter(b)

		if len(r.PostForm) > 0 {
			for key, values := range r.PostForm {
				for _, value := range values {
					if err := w.WriteField(key, value); err != nil {
						return err
					}
				}
			}
		}

		for _, file := range files {
			// custom header since CreateFormFile uses application/octet-stream by default
			var fw io.Writer
			var err error
			if file.ContentType != "" {
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition",
					fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
						escapeQuotes(file.Key), escapeQuotes(file.FileName)))
				h.Set("Content-Type", file.ContentType)
				fw, err = w.CreatePart(h)
				if err != nil {
					return err
				}
			} else {
				fw, err = w.CreateFormFile(file.Key, file.FileName)
				if err != nil {
					return err
				}
			}
			_, err = io.Copy(fw, bytes.NewBuffer(file.Contents))
			if err != nil {
				return err
			}
		}
		r.Header.Set(HeaderContentType, w.FormDataContentType())
		if err := w.Close(); err != nil {
			return err
		}

		bb := b.Bytes()
		r.Body = io.NopCloser(bytes.NewReader(bb))
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bb)), nil
		}
		r.ContentLength = int64(len(bb))
		return nil
	}
}

// OptJSONBody sets the post body on the request.
func OptJSONBody(obj interface{}) RequestOption {
	return func(r *http.Request) error {
		contents, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		r.Body = io.NopCloser(bytes.NewReader(contents))
		r.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(contents)
			return io.NopCloser(r), nil
		}
		r.ContentLength = int64(len(contents))
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationJSON)
		return nil
	}
}

// OptXMLBody sets the post body on the request.
func OptXMLBody(obj interface{}) RequestOption {
	return func(r *http.Request) error {
		contents, err := xml.Marshal(obj)
		if err != nil {
			return err
		}
		r.Body = io.NopCloser(bytes.NewBuffer(contents))
		r.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(contents)
			return io.NopCloser(r), nil
		}
		r.ContentLength = int64(len(contents))
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set(HeaderContentType, ContentTypeApplicationXML)
		return nil
	}
}

// OptHTTPClientTrace sets the http trace on the outgoing request.
func OptHTTPClientTrace(ht *HTTPTrace) RequestOption {
	return func(r *http.Request) error {
		*r = *WithClientHTTPTrace(r, ht)
		return nil
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
