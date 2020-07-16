package spiffeutil_test

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/spiffeutil"
)

func TestParse(t *testing.T) {
	it := assert.New(t)

	type FailureCase struct {
		URI     string
		Message string
	}
	failures := []FailureCase{
		FailureCase{URI: "https://web.invalid", Message: "Does not match protocol: \"https://web.invalid\""},
		FailureCase{URI: "spiffe://only.local", Message: "Missing workload identifier: \"spiffe://only.local\""},
		FailureCase{URI: "spiffe://only.local/", Message: "Missing workload identifier: \"spiffe://only.local/\""},
	}
	for _, fc := range failures {
		trustDomain, workloadID, err := spiffeutil.Parse(fc.URI)
		it.Empty(trustDomain)
		it.Empty(workloadID)
		it.True(ex.Is(err, spiffeutil.ErrInvalidURI))
		asEx, ok := err.(*ex.Ex)
		it.True(ok)
		it.Equal(fc.Message, asEx.Message)
	}

	// Success.
	trustDomain, workloadID, err := spiffeutil.Parse("spiffe://cluster.local/ns/blend/sa/quasar")
	it.Equal("cluster.local", trustDomain)
	it.Equal("ns/blend/sa/quasar", workloadID)
	it.Nil(err)
}
