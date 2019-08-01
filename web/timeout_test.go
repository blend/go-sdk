package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimeout(t *testing.T) {
	assert := assert.New(t)

	app := MustNew(
		OptBindAddr(DefaultMockBindAddr),
		OptUse(WithTimeout(1*time.Millisecond)),
	)

	var didFinish bool
	app.GET("/panic", func(_ *Ctx) Result {
		panic("test")
	})
	app.GET("/long", func(_ *Ctx) Result {
		time.Sleep(4 * time.Millisecond)
		didFinish = true
		return NoContent
	})
	app.GET("/short", func(_ *Ctx) Result {
		didFinish = true
		return NoContent
	})

	go app.Start()
	defer app.Stop()
	<-app.NotifyStarted()

	_, err := http.Get("http://" + app.Listener.Addr().String() + "/panic")
	assert.Nil(err)
	assert.False(didFinish)

	_, err = http.Get("http://" + app.Listener.Addr().String() + "/long")
	assert.Nil(err)
	assert.False(didFinish)

	_, err = http.Get("http://" + app.Listener.Addr().String() + "/short")
	assert.Nil(err)
	assert.True(didFinish)
}
