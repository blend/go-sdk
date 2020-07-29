package envoyutil_test

import (
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/collections"
	"github.com/blend/go-sdk/spiffeutil"
	"github.com/blend/go-sdk/uuid"

	"github.com/blend/go-sdk/envoyutil"
)

func TestOptAllowedTrustDomains(t *testing.T) {
	assert := sdkAssert.New(t)

	cip := &envoyutil.ClientIdentityProcessor{
		AllowedTrustDomains: []string{"x.invalid"},
	}
	opt := envoyutil.OptAllowedTrustDomains("y.invalid")
	opt(cip)

	expected := &envoyutil.ClientIdentityProcessor{
		AllowedTrustDomains: []string{"x.invalid", "y.invalid"},
	}
	assert.Equal(expected, cip)
}

func TestOptDeniedTrustDomains(t *testing.T) {
	assert := sdkAssert.New(t)

	cip := &envoyutil.ClientIdentityProcessor{
		DeniedTrustDomains: []string{"y.invalid"},
	}
	opt := envoyutil.OptDeniedTrustDomains("z.invalid")
	opt(cip)

	expected := &envoyutil.ClientIdentityProcessor{
		DeniedTrustDomains: []string{"y.invalid", "z.invalid"},
	}
	assert.Equal(expected, cip)
}

func TestOptAllowedClientIdentities(t *testing.T) {
	assert := sdkAssert.New(t)

	cip := &envoyutil.ClientIdentityProcessor{
		AllowedClientIdentities: collections.NewSetOfString("x.invalid", "y.invalid"),
	}
	opt := envoyutil.OptAllowedClientIdentities("y.invalid", "z.invalid")
	opt(cip)

	expected := &envoyutil.ClientIdentityProcessor{
		AllowedClientIdentities: collections.NewSetOfString("x.invalid", "y.invalid", "z.invalid"),
	}
	assert.Equal(expected, cip)
}

func TestOptDeniedClientIdentities(t *testing.T) {
	assert := sdkAssert.New(t)

	cip := &envoyutil.ClientIdentityProcessor{
		DeniedClientIdentities: collections.NewSetOfString("x.invalid", "y.invalid"),
	}
	opt := envoyutil.OptDeniedClientIdentities("y.invalid", "z.invalid")
	opt(cip)

	expected := &envoyutil.ClientIdentityProcessor{
		DeniedClientIdentities: collections.NewSetOfString("x.invalid", "y.invalid", "z.invalid"),
	}
	assert.Equal(expected, cip)
}

func TestOptFormatClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	cip := &envoyutil.ClientIdentityProcessor{
		FormatClientIdentity: envoyutil.KubernetesClientIdentityFormatter,
	}
	sentinel := uuid.V4().ToFullString()
	var fn envoyutil.ClientIdentityFormatter = func(_ envoyutil.XFCCElement, _ *spiffeutil.ParsedURI) (string, error) {
		return sentinel, nil
	}
	opt := envoyutil.OptFormatClientIdentity(fn)
	opt(cip)

	// Can't compare functions for equality, see https://github.com/blend/go-sdk/issues/167
	// so we make sure our function is as expected.
	s, err := cip.FormatClientIdentity(envoyutil.XFCCElement{}, nil)
	assert.Equal(sentinel, s)
	assert.Nil(err)
}
