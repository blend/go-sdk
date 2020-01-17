package reverseproxy

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptProxyTransformRequest(t *testing.T) {
	it := assert.New(t)

	var calledReq *http.Request
	tr := func(req *http.Request) error {
		calledReq = req
		return nil
	}
	p := NewProxy(OptProxyTransformRequest(tr))

	// Need to special case function equality.
	it.Equal(reflect.ValueOf(tr).Pointer(), reflect.ValueOf(p.TransformRequest).Pointer())
	it.Nil(calledReq)
}
