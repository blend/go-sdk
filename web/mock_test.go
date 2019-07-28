package web

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMock(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.GET("/", func(_ *Ctx) Result { return NoContent })

	res, err := Mock(app, &http.Request{Method: "GET", URL: &url.URL{Scheme: SchemeHTTP, Path: "/"}}).DiscardWithResponse()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)

	assert.True(app.IsStopped())
}

func TestMockGet(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.GET("/", func(_ *Ctx) Result { return NoContent })

	res, err := MockGet(app, "/").DiscardWithResponse()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)

	assert.True(app.IsStopped())
}
