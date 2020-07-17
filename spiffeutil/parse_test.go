package spiffeutil_test

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/spiffeutil"
)

func TestParse(t *testing.T) {
	it := assert.New(t)

	type failureCase struct {
		URI     string
		Message string
	}
	failures := []failureCase{
		{URI: "https://web.invalid", Message: "Does not match protocol: \"https://web.invalid\""},
		{URI: "spiffe://only.local", Message: "Missing workload identifier: \"spiffe://only.local\""},
		{URI: "spiffe://only.local/", Message: "Missing workload identifier: \"spiffe://only.local/\""},
	}
	for _, fc := range failures {
		pu, err := spiffeutil.Parse(fc.URI)
		it.Nil(pu)
		it.True(ex.Is(err, spiffeutil.ErrInvalidURI))
		asEx, ok := err.(*ex.Ex)
		it.True(ok)
		it.Equal(fc.Message, asEx.Message)
	}

	// Success.
	pu, err := spiffeutil.Parse("spiffe://cluster.local/ns/blend/sa/quasar")
	expected := &spiffeutil.ParsedURI{TrustDomain: "cluster.local", WorkloadID: "ns/blend/sa/quasar"}
	it.Equal(expected, pu)
	it.Nil(err)
}
