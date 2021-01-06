package web

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestDefaultProviderMiddlewares(t *testing.T) {
	assert := assert.New(t)

	r := applyMiddleware(JSONProviderAsDefault)
	_, ok := r.DefaultProvider.(JSONResultProvider)
	assert.True(ok)

	r = applyMiddleware(ViewProviderAsDefault)
	_, ok = r.DefaultProvider.(*ViewCache)
	assert.True(ok)

	r = applyMiddleware(XMLProviderAsDefault)
	_, ok = r.DefaultProvider.(XMLResultProvider)
	assert.True(ok)

	r = applyMiddleware(TextProviderAsDefault)
	_, ok = r.DefaultProvider.(TextResultProvider)
	assert.True(ok)
}

func applyMiddleware(middleware Middleware) (output *Ctx) {
	middleware(func(ctx *Ctx) Result {
		output = ctx
		return NoContent
	})(NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/")))
	return
}
