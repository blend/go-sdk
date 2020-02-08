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

	var didShortFinish, didLongFinish bool
	app.GET("/panic", func(_ *Ctx) Result {
		panic("test")
	})
	app.GET("/long", func(_ *Ctx) Result {
		time.Sleep(4 * time.Millisecond)
		didLongFinish = true
		return NoContent
	})
	app.GET("/short", func(_ *Ctx) Result {
		didShortFinish = true
		return NoContent
	})

	go app.Start()
	defer app.Stop()
	<-app.NotifyStarted()

	_, err := http.Get("http://" + app.Listener.Addr().String() + "/panic")
	assert.Nil(err)

	_, err = http.Get("http://" + app.Listener.Addr().String() + "/long")
	assert.Nil(err)
	assert.False(didLongFinish)

	_, err = http.Get("http://" + app.Listener.Addr().String() + "/short")
	assert.Nil(err)
	assert.True(didShortFinish)
}
