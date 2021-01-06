package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSSkipVerify(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTLSSkipVerify(true))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
}
