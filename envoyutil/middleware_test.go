package envoyutil_test

import (
	"encoding/json"
	"net/http"
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
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

func TestClientIdentityRequired(t *testing.T) {
	assert := sdkAssert.New(t)

	app := web.MustNew()
	var capturedContext *web.Ctx
	cip := envoyutil.ClientIdentityFromSPIFFE(
		envoyutil.OptDeniedIdentities("gw.blend"),
	)
	verifier := envoyutil.ServerIdentityFromSPIFFE(
		envoyutil.OptAllowedIdentities("idea.blend"),
	)
	app.GET(
		"/",
		func(ctx *web.Ctx) web.Result {
			capturedContext = ctx
			return web.JSON.OK()
		},
		envoyutil.ClientIdentityRequired(cip, verifier),
		web.JSONProviderAsDefault,
	)

	body, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on missing header")
	assert.Nil(capturedContext)
	var expected error = &envoyutil.XFCCValidationError{Class: envoyutil.ErrMissingXFCC}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc := `""`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on empty header")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{Class: envoyutil.ErrInvalidXFCC, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = "something=bad"
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on malformed header")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{Class: envoyutil.ErrInvalidXFCC, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = "By=spiffe://cluster.local/ns/blend/sa/idea;URI=spiffe://cluster.local/ns/blend/sa/should-end/sa/extra"
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on unexpected SPIFFE format in `URI`")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{Class: envoyutil.ErrInvalidClientIdentity, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/idea;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client"`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on missing client identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCValidationError{Class: envoyutil.ErrInvalidClientIdentity, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/idea;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/gw`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on denied client identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCValidationError{
		Class: envoyutil.ErrDeniedClientIdentity,
		XFCC:  xfcc,
		// NOTE: This should really be a `map[string]string`. We use a `map[string]interface{}`
		//       so that the comparison in `invalidXFCCJSONEqual()` passes.
		Metadata: map[string]interface{}{"clientIdentity": "gw.blend"},
	}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/idea;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/twtr`
	meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Success on valid header")
	assert.NotNil(capturedContext)
	assert.Equal("twtr.blend", envoyutil.GetClientIdentity(capturedContext))
	capturedContext = nil

	xfcc = `By=mailto:John.Doe@example.com;URI=spiffe://cluster.local/ns/blend/sa/peas`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on invalid server identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{
		Class: envoyutil.ErrInvalidServerIdentity,
		XFCC:  xfcc,
	}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/outside;URI=spiffe://cluster.local/ns/blend/sa/peas`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on wrong server identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCValidationError{
		Class: envoyutil.ErrDeniedServerIdentity,
		XFCC:  xfcc,
		// NOTE: This should really be a `map[string]string`. We use a `map[string]interface{}`
		//       so that the comparison in `invalidXFCCJSONEqual()` passes.
		Metadata: map[string]interface{}{"serverIdentity": "outside.blend"},
	}
	invalidXFCCJSONEqual(assert, expected, body)

	// Unrecoverable error: here we simulate `envoyutil` user error by using
	// `nil` for `cip`.
	app = web.MustNew()
	app.GET(
		"/",
		func(ctx *web.Ctx) web.Result {
			return web.JSON.OK()
		},
		envoyutil.ClientIdentityRequired(nil),
		web.JSONProviderAsDefault,
	)
	body, meta, err = web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, meta.StatusCode, "Fail on unrecoverable")
	assert.Equal("\"Internal Server Error\"\n", string(body))
}

func invalidXFCCJSONEqual(assert *sdkAssert.Assertions, expected error, actual []byte) {
	switch expected.(type) {
	case *envoyutil.XFCCExtractionError:
		unmarshaledActual := &envoyutil.XFCCExtractionError{}
		err := json.Unmarshal(actual, unmarshaledActual)
		assert.Nil(err)
		assert.Equal(expected, unmarshaledActual)
	case *envoyutil.XFCCValidationError:
		unmarshaledActual := &envoyutil.XFCCValidationError{}
		err := json.Unmarshal(actual, unmarshaledActual)
		assert.Nil(err)
		assert.Equal(expected, unmarshaledActual)
	default:
		assert.FailNow()
	}
}
