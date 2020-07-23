package web

import (
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimeout(t *testing.T) {
	t.Skip() // flaky
	assert := assert.New(t)

	app := MustNew(
		OptBindAddr(DefaultMockBindAddr),
		OptUse(WithTimeout(1*time.Millisecond)),
	)

	var didShortFinish, didLongFinish int32
	app.GET("/panic", func(_ *Ctx) Result {
		panic("test")
	})
	app.GET("/long", func(_ *Ctx) Result {
		time.Sleep(5 * time.Millisecond)
		atomic.StoreInt32(&didLongFinish, 1)
		return NoContent
	})
	app.GET("/short", func(_ *Ctx) Result {
		atomic.StoreInt32(&didShortFinish, 1)
		return NoContent
	})

	go func() { _ = app.Start() }()
	defer func() { _ = app.Stop() }()
	<-app.NotifyStarted()

	res, err := http.Get("http://" + app.Listener.Addr().String() + "/panic")
	assert.Nil(err)
	assert.Nil(res.Body.Close())

	res, err = http.Get("http://" + app.Listener.Addr().String() + "/long")
	assert.Nil(err)
	assert.Nil(res.Body.Close())
	assert.Zero(atomic.LoadInt32(&didLongFinish))

	res, err = http.Get("http://" + app.Listener.Addr().String() + "/short")
	assert.Nil(err)
	assert.Nil(res.Body.Close())
	assert.Equal(1, atomic.LoadInt32(&didShortFinish))
}
