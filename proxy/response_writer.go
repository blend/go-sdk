package proxy

import (
	"net/http"
)

// NewResponseWriter creates a new uncompressed response writer.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		innerResponse: w,
	}
}

// ResponseWriter a better response writer
type ResponseWriter struct {
	innerResponse http.ResponseWriter
	http.Hijacker
	contentLength int
	statusCode    int
}

// Write writes the data to the response.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	written, err := rw.innerResponse.Write(b)
	rw.contentLength += written
	return written, err
}

// Header accesses the response header collection.
func (rw *ResponseWriter) Header() http.Header {
	return rw.innerResponse.Header()
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

// StatusCode returns the status code.
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// ContentLength returns the content length
func (rw *ResponseWriter) ContentLength() int {
	return rw.contentLength
}

// Close disposes of the response writer.
func (rw *ResponseWriter) Close() error {
	return nil
}
