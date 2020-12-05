package webutil

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

// NOTE: Ensure `mockResponseWriterWithHijack` satisfies `http.Hijacker`.
var (
	_ http.Hijacker = mockResponseWriterWithHijack{}
)

type mockResponseWriter struct {
	Headers    http.Header
	StatusCode int
	Output     io.Writer
}

// Header returns the response headers.
func (mrw mockResponseWriter) Header() http.Header {
	return mrw.Headers
}

// WriteHeader writes the status code.
func (mrw mockResponseWriter) WriteHeader(code int) {
	mrw.StatusCode = code
}

// Write writes data.
func (mrw mockResponseWriter) Write(contents []byte) (int, error) {
	return mrw.Output.Write(contents)
}

type mockResponseWriterWithHijack struct {
	mockResponseWriter
	HijackError error
}

func (mrwwh mockResponseWriterWithHijack) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, mrwwh.HijackError
}

func TestNewResponseWriter(t *testing.T) {
	assert := assert.New(t)

	var ir http.ResponseWriter
	ir = mockResponseWriter{}
	rw1 := &ResponseWriter{innerResponse: ir}
	rw2 := NewResponseWriter(rw1)
	assert.Equal(rw1, rw2)
}

func TestResponseWriter(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	rw := NewResponseWriter(mockResponseWriter{Output: output, Headers: http.Header{}})

	rw.Header().Set("foo", "bar")
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("this is a test"))
	assert.Nil(err)

	assert.Equal(http.StatusOK, rw.StatusCode())
	assert.Equal("this is a test", output.String())
}

func TestResponseWriterHijack(t *testing.T) {
	assert := assert.New(t)

	// **Does not** support Hijack
	var ir http.ResponseWriter
	ir = mockResponseWriter{}
	rw := ResponseWriter{innerResponse: ir}
	nc, buf, err := rw.Hijack()
	assert.Nil(nc)
	assert.Nil(buf)
	assert.True(ex.Is(ErrHijackerUnsupported, err))

	// **Does** support Hijack
	hijackError := ex.New("local failure")
	ir = mockResponseWriterWithHijack{HijackError: hijackError}
	rw = ResponseWriter{innerResponse: ir}
	nc, buf, err = rw.Hijack()
	assert.Nil(nc)
	assert.Nil(buf)
	assert.ReferenceEqual(hijackError, err)
}

func TestResponseWriterInnerResponse(t *testing.T) {
	assert := assert.New(t)

	var ir http.ResponseWriter
	ir = mockResponseWriter{}
	rw := ResponseWriter{innerResponse: ir}
	assert.Equal(ir, rw.InnerResponse())
}

func TestResponseWriterFlush(t *testing.T) {
	assert := assert.New(t)

	rw := ResponseWriter{}
	rw.Flush()
	assert.True(true)
}
