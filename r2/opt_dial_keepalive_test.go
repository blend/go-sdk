package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptDialKeepAlive(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptDial(OptDialKeepAlive(time.Second)))
	assert.NotNil(r.Client.Transport.(*http.Transport).DialContext)
}
