package envoyutil_test

import (
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/web"

	"github.com/blend/go-sdk/envoyutil"
)

func TestGetClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	ctx := web.MockCtx("GET", "/")
	assert.Empty(envoyutil.GetClientIdentity(ctx))

	ctx.WithStateValue(envoyutil.StateKeyClientIdentity, nil)
	assert.Empty(envoyutil.GetClientIdentity(ctx))

	ctx.WithStateValue(envoyutil.StateKeyClientIdentity, 42)
	assert.Empty(envoyutil.GetClientIdentity(ctx))

	wi := "hello.world"
	ctx.WithStateValue(envoyutil.StateKeyClientIdentity, wi)
	assert.NotNil(envoyutil.GetClientIdentity(ctx))
	assert.Equal(wi, envoyutil.GetClientIdentity(ctx))
}
