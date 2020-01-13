package webutil

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"github.com/blend/go-sdk/ex"
)

var (
	_ (ResponseWrapper)     = (*ResponseWriter)(nil)
	_ (http.ResponseWriter) = (*ResponseWriter)(nil)
	_ (http.Flusher)        = (*ResponseWriter)(nil)
)

// NewResponseWriter creates a new response writer.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	if typed, ok := w.(ResponseWrapper); ok {
		return &ResponseWriter{
			innerResponse: typed.InnerResponse(),
		}
	}
	return &ResponseWriter{
		innerResponse: w,
	}
}

// ResponseWriter a better response writer
type ResponseWriter struct {
	innerResponse http.ResponseWriter
	statusCode    int
	contentLength int
}

// Write writes the data to the response.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	bytesWritten, err := rw.innerResponse.Write(b)
	rw.contentLength = rw.contentLength + bytesWritten
	return bytesWritten, err
}

// Header accesses the response header collection.
func (rw *ResponseWriter) Header() http.Header {
	return rw.innerResponse.Header()
}

// Hijack wraps response writer's Hijack function.
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.innerResponse.(http.Hijacker)
	if !ok {
		return nil, nil, ex.New(fmt.Errorf("ResponseWriter doesn't support Hijacker interface"))
	}
	return hijacker.Hijack()
}

// WriteHeader is actually a terrible name and this writes the status code.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.innerResponse.WriteHeader(code)
}

// InnerResponse returns the backing writer.
func (rw *ResponseWriter) InnerResponse() http.ResponseWriter {
	return rw.innerResponse
}

// Flush is a no op on raw response writers.
func (rw *ResponseWriter) Flush() {}

// StatusCode returns the status code.
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// ContentLength returns the content length
func (rw *ResponseWriter) ContentLength() int {
	return rw.contentLength
}
