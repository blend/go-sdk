package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

type mockTracer struct{}

func (mt mockTracer) Start(_ *http.Request) TraceFinisher {
	return mockTraceFinisher{}
}

type mockTraceFinisher struct{}

func (mtf mockTraceFinisher) Finish(_ *http.Request, _ *http.Response, _ time.Time, _ error) {}

func TestOptTracer(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptTracer(mockTracer{}))
	assert.NotNil(r.Tracer)
}
