package web

import (
	"context"
	"net/http"
	"testing"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/stringutil"
)

func TestSessionAware(t *testing.T) {
	assert := assert.New(t)

	sessionID := NewSessionID()

	var didExecuteHandler bool
	var sessionWasSet bool

	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"})

	app.GET("/", func(r *Ctx) Result {
		didExecuteHandler = true
		sessionWasSet = r.Session != nil
		return Text.Result("COOL")
	}, SessionAware)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(ContentTypeText, meta.Header.Get(HeaderContentType))
	assert.True(didExecuteHandler, "we should have triggered the hander")
	assert.True(sessionWasSet, "the session should have been set by the middleware")

	unsetMeta, err := MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, unsetMeta.StatusCode)
	assert.False(sessionWasSet)
}

func TestSessionRequired(t *testing.T) {
	assert := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"})

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		return Text.Result("COOL")
	}, SessionRequired)

	unsetMeta, err := MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, unsetMeta.StatusCode)
	assert.False(sessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(sessionWasSet)
}

func TestSessionRequiredCustomParamName(t *testing.T) {
	assert := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"})
	app.Auth.CookieDefaults.Name = "web_auth"

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		return Text.Result("COOL")
	}, SessionRequired)

	unsetMeta, err := MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, unsetMeta.StatusCode)
	assert.False(sessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(sessionWasSet)

	meta, err = MockGet(app, "/", r2.OptCookieValue(DefaultCookieName, sessionID)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode)
	assert.True(sessionWasSet)
}

func TestSessionMiddleware(t *testing.T) {
	assert := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()), OptBindAddr(DefaultMockBindAddr))
	app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "bailey"})

	go app.Start()
	<-app.NotifyStarted()
	defer app.Stop()

	var calledCustom bool
	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		return Text.Result("COOL")
	}, SessionMiddleware(func(_ *Ctx) Result {
		calledCustom = true
		return NoContent
	}))

	unsetMeta, err := MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, unsetMeta.StatusCode)
	assert.False(sessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(sessionWasSet)
	assert.True(calledCustom)
}
