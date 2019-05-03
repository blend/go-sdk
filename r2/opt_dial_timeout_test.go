package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptDialTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptDial(OptDialTimeout(time.Second)))
	assert.NotNil(r.Client.Transport.(*http.Transport).Dial)
}
