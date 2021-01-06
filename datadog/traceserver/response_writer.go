package traceserver

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

// ResponseWriter a better response writer
type ResponseWriter struct {
	http.ResponseWriter

	StatusCode    int
	ContentLength int
}

// Write writes the data to the response.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	bytesWritten, err := rw.ResponseWriter.Write(b)
	rw.ContentLength += bytesWritten
	return bytesWritten, err
}

// Hijack wraps response writer's Hijack function.
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer hijack; wrapped response writer doesn't support the hijacker interface")
	}
	return hijacker.Hijack()
}

// WriteHeader writes the status code (it is a somewhat poorly chosen method name from the standard library).
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Flush calls flush on the inner response writer if it is supported.
func (rw *ResponseWriter) Flush() {
	if typed, ok := rw.ResponseWriter.(http.Flusher); ok {
		typed.Flush()
	}
}

// Close calls close on the inner response if it supports it.
func (rw *ResponseWriter) Close() error {
	if typed, ok := rw.ResponseWriter.(io.Closer); ok {
		return typed.Close()
	}
	return nil
}
