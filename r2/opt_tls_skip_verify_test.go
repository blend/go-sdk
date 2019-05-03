package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSSkipVerify(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptTLSSkipVerify(true))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
}
