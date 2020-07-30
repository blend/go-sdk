package envoyutil_test

import (
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/spiffeutil"
	"github.com/blend/go-sdk/uuid"

	"github.com/blend/go-sdk/envoyutil"
)

func TestOptIdentityType(t *testing.T) {
	assert := sdkAssert.New(t)

	ip := &envoyutil.IdentityProcessor{}
	opt := envoyutil.OptIdentityType(envoyutil.ServerIdentity)
	opt(ip)

	expected := &envoyutil.IdentityProcessor{
		Type: envoyutil.ServerIdentity,
	}
	assert.Equal(expected, ip)
}

func TestOptFormatIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	ip := &envoyutil.IdentityProcessor{
		FormatIdentity: makeMockFormatter("not-here"),
	}
	sentinel := uuid.V4().ToFullString()
	formatter := makeMockFormatter(sentinel)
	opt := envoyutil.OptFormatIdentity(formatter)
	opt(ip)

	// Can't compare functions for equality, see https://github.com/blend/go-sdk/issues/167
	// so we make sure our function is as expected.
	s, err := ip.FormatIdentity(envoyutil.XFCCElement{}, nil)
	assert.Equal(sentinel, s)
	assert.Nil(err)
}

func TestIdentityProcessorKubernetesIdentityFormatter(t *testing.T) {
	assert := sdkAssert.New(t)

	xfcc := envoyutil.XFCCElement{By: "anything", URI: "goes"}

	// Valid identity.
	ip := &envoyutil.IdentityProcessor{}
	pu := &spiffeutil.ParsedURI{WorkloadID: "ns/packets/sa/ketchup"}
	identity, err := ip.KubernetesIdentityFormatter(xfcc, pu)
	assert.Equal("ketchup.packets", identity)
	assert.Nil(err)

	// Invalid client identity.
	pu = &spiffeutil.ParsedURI{WorkloadID: "not-k8s"}
	clientIdentity, err := ip.KubernetesIdentityFormatter(xfcc, pu)
	assert.Equal("", clientIdentity)
	assert.True(envoyutil.IsExtractionError(err))
	expected := &envoyutil.XFCCExtractionError{
		Class: envoyutil.ErrInvalidClientIdentity,
		XFCC:  xfcc.String(),
	}
	assert.Equal(expected, err)

	// Invalid server identity.
	ip = &envoyutil.IdentityProcessor{Type: envoyutil.ServerIdentity}
	serverIdentity, err := ip.KubernetesIdentityFormatter(xfcc, pu)
	assert.Equal("", serverIdentity)
	assert.True(envoyutil.IsExtractionError(err))
	expected = &envoyutil.XFCCExtractionError{
		Class: envoyutil.ErrInvalidServerIdentity,
		XFCC:  xfcc.String(),
	}
	assert.Equal(expected, err)
}

func makeMockFormatter(identity string) envoyutil.IdentityFormatter {
	return func(_ envoyutil.XFCCElement, _ *spiffeutil.ParsedURI) (string, error) {
		return identity, nil
	}
}
