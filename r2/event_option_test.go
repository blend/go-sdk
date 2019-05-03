package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestEventOptions(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent("http.client.request")

	assert.Equal(e.Flag, "http.client.request")
	OptEventFlag("http.client.response")(e)
	assert.Equal(e.Flag, "http.client.response")

	t0 := time.Date(2019, 05, 02, 12, 13, 14, 15, time.UTC)
	assert.NotEqual(t0, e.Timestamp)
	OptEventCompleted(t0)(e)
	assert.Equal(t0, e.Timestamp)

	t1 := time.Date(2019, 05, 02, 12, 13, 14, 55, time.UTC)
	assert.NotEqual(t1, e.Started)
	OptEventStarted(t1)(e)
	assert.Equal(t1, e.Started)

	assert.Nil(e.Request)
	OptEventRequest(&http.Request{RequestURI: "abcdef"})(e)
	assert.NotNil(e.Request)
	assert.Equal("abcdef", e.Request.RequestURI)

	assert.Nil(e.Response)
	OptEventResponse(&http.Response{Proto: "not-http"})(e)
	assert.NotNil(e.Response)
	assert.Equal("not-http", e.Response.Proto)

	assert.Nil(e.Body)
	OptEventBody([]byte(`bailey`))(e)
	assert.NotNil(e.Body)
	assert.Equal("bailey", string(e.Body))
}
