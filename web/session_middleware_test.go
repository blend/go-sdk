/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/webutil"
)

func Test_SessionAware(t *testing.T) {
	its := assert.New(t)

	sessionID := NewSessionID()

	var didExecuteHandler bool
	var sessionWasSet bool
	var sessionExpirationWasChanged bool
	var contextSessionWasSet bool

	now := time.Time{}
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.SessionTimeoutProvider = func(s *Session) time.Time {
		return now.Add(time.Hour)
	}
	err := app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"})
	its.Nil(err)

	_, _, err = app.Auth.VerifySession(MockCtx(http.MethodGet, "/"))
	its.Nil(err)

	app.GET("/", func(r *Ctx) Result {
		didExecuteHandler = true
		sessionWasSet = r.Session != nil
		if r.Session != nil {
			sessionExpirationWasChanged = r.Session.ExpiresUTC != now
		}
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionAware)

	unsetMeta, err := MockGet(app, "/").Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, unsetMeta.StatusCode)
	its.False(sessionWasSet)
	its.False(sessionExpirationWasChanged)
	its.False(contextSessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.Equal(webutil.ContentTypeText, meta.Header.Get(webutil.HeaderContentType))
	its.True(didExecuteHandler, "we should have triggered the hander")
	its.True(sessionWasSet, "the session should have been set by the middleware")
	its.True(sessionExpirationWasChanged, "the session should have had its expiration updated")
	its.True(contextSessionWasSet, "the context session should have been set by the middleware")
}

func Test_SessionAware_errSessionInvalid(t *testing.T) {
	its := assert.New(t)

	sessionID := NewSessionID()

	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return ErrSessionIDEmpty
	}

	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	var didExecuteHandler bool
	var sessionWasSet bool
	var contextSessionWasSet bool
	app.GET("/", func(r *Ctx) Result {
		didExecuteHandler = true
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionAware)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.Equal(webutil.ContentTypeText, meta.Header.Get(webutil.HeaderContentType))
	its.True(didExecuteHandler, "we should have triggered the hander")
	its.False(sessionWasSet, "the session should not have been set by the middleware")
	its.False(contextSessionWasSet, "the context session should not have been set by the middleware")
}

func Test_SessionAware_error(t *testing.T) {
	its := assert.New(t)

	sessionID := NewSessionID()

	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return fmt.Errorf("this is just a test")
	}

	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	var didExecuteHandler bool
	var sessionWasSet bool
	var contextSessionWasSet bool
	app.GET("/", func(r *Ctx) Result {
		didExecuteHandler = true
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionAware)

	meta, err := MockGet(app, "/",
		r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID),
	).Discard()
	its.Nil(err)
	its.Equal(http.StatusInternalServerError, meta.StatusCode)
	its.False(didExecuteHandler, "we should have triggered the hander")
	its.False(sessionWasSet, "the session should not have been set by the middleware")
	its.False(contextSessionWasSet, "the context session should not have been set by the middleware")
}

func Test_SessionAwareForLogout(t *testing.T) {
	its := assert.New(t)

	sessionID := NewSessionID()

	var didExecuteHandler bool
	var sessionWasSet bool
	var sessionExpirationWasChanged bool
	var contextSessionWasSet bool

	now := time.Time{}
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.SessionTimeoutProvider = func(s *Session) time.Time {
		return now.Add(time.Hour)
	}

	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	app.GET("/", func(r *Ctx) Result {
		didExecuteHandler = true
		sessionWasSet = r.Session != nil
		if r.Session != nil {
			sessionExpirationWasChanged = r.Session.ExpiresUTC != now
		}
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionAwareForLogout)

	unsetMeta, err := MockGet(app, "/").Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, unsetMeta.StatusCode)
	its.False(sessionWasSet)
	its.False(sessionExpirationWasChanged)
	its.False(contextSessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.Equal(webutil.ContentTypeText, meta.Header.Get(webutil.HeaderContentType))
	its.True(didExecuteHandler, "we should have triggered the hander")
	its.True(sessionWasSet, "the session should have been set by the middleware")
	its.False(sessionExpirationWasChanged, "we should _not_ have updated the session expiry")
	its.True(contextSessionWasSet, "the context session should have been set by the middleware")
}

func Test_SessionAwareForLogout_error(t *testing.T) {
	its := assert.New(t)

	sessionID := NewSessionID()

	var sessionWasSet bool
	var sessionExpirationWasChanged bool
	var contextSessionWasSet bool

	now := time.Time{}
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.SessionTimeoutProvider = func(s *Session) time.Time {
		return now.Add(time.Hour)
	}
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return fmt.Errorf("this is just a test")
	}

	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		if r.Session != nil {
			sessionExpirationWasChanged = r.Session.ExpiresUTC != now
		}
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionAwareForLogout)

	unsetMeta, err := MockGet(app, "/").Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, unsetMeta.StatusCode)
	its.False(sessionWasSet)
	its.False(sessionExpirationWasChanged)
	its.False(contextSessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusInternalServerError, meta.StatusCode)
	its.False(sessionWasSet, "the session should not have been set by the middleware")
	its.False(sessionExpirationWasChanged, "we should not have updated the session expiry")
	its.False(contextSessionWasSet, "the context session should not have been set by the middleware")
}

func Test_SessionRequired(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()))
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionRequired)

	unsetMeta, err := MockGet(app, "/").Discard()
	its.Nil(err)
	its.Equal(http.StatusUnauthorized, unsetMeta.StatusCode)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(sessionWasSet)
	its.True(contextSessionWasSet)
}

func Test_SessionRequired_errSessionInvalid(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var didCallHandler bool
	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return ErrSessionIDEmpty
	}
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	app.GET("/", func(r *Ctx) Result {
		didCallHandler = true
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionRequired)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusUnauthorized, meta.StatusCode)
	its.False(didCallHandler)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)
}

func Test_SessionRequired_error(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var didCallHandler bool
	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return fmt.Errorf("this is just a test")
	}
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	app.GET("/", func(r *Ctx) Result {
		didCallHandler = true
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionRequired)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusInternalServerError, meta.StatusCode)
	its.False(didCallHandler)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)
}

func Test_SessionRequired_customParamName(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()))
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))
	app.Auth.CookieDefaults.Name = "web_auth"

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionRequired)

	unsetMeta, err := MockGet(app, "/").Discard()
	its.Nil(err)
	its.Equal(http.StatusUnauthorized, unsetMeta.StatusCode)
	its.False(sessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(sessionWasSet)
	its.True(contextSessionWasSet)

	meta, err = MockGet(app, "/", r2.OptCookieValue(DefaultCookieName, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusUnauthorized, meta.StatusCode)
}

func Test_SessionMiddleware(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()), OptBindAddr(DefaultMockBindAddr))
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	go func() { _ = app.Start() }()
	<-app.NotifyStarted()
	defer func() { _ = app.Stop() }()

	var calledCustom bool
	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionMiddleware(func(_ *Ctx) Result {
		calledCustom = true
		return NoContent
	}))

	unsetMeta, err := MockGet(app, "/").Discard()
	its.Nil(err)
	its.Equal(http.StatusNoContent, unsetMeta.StatusCode)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(sessionWasSet)
	its.True(calledCustom)
}

func Test_SessionMiddleware_errSessionInvalid(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()), OptBindAddr(DefaultMockBindAddr))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return ErrSessionIDEmpty
	}
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	go func() { _ = app.Start() }()
	<-app.NotifyStarted()
	defer func() { _ = app.Stop() }()

	var calledCustom bool
	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionMiddleware(func(_ *Ctx) Result {
		calledCustom = true
		return NoContent
	}))

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusNoContent, meta.StatusCode)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)
	its.True(calledCustom)
}

func Test_SessionMiddleware_errSessionInvalid_unsetCustom(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()), OptBindAddr(DefaultMockBindAddr))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return ErrSessionIDEmpty
	}
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	go func() { _ = app.Start() }()
	<-app.NotifyStarted()
	defer func() { _ = app.Stop() }()

	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionMiddleware(nil))

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusUnauthorized, meta.StatusCode)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)
}

func Test_SessionMiddleware_error(t *testing.T) {
	its := assert.New(t)

	sessionID := stringutil.Random(stringutil.LettersAndNumbers, 64)

	var sessionWasSet bool
	var contextSessionWasSet bool
	app := MustNew(OptAuth(NewLocalAuthManager()), OptBindAddr(DefaultMockBindAddr))
	app.Auth.ValidateHandler = func(_ context.Context, _ *Session) error {
		return fmt.Errorf("this is just a test")
	}
	its.Nil(app.Auth.PersistHandler(context.TODO(), &Session{SessionID: sessionID, UserID: "example-string"}))

	go func() { _ = app.Start() }()
	<-app.NotifyStarted()
	defer func() { _ = app.Stop() }()

	var calledCustom bool
	app.GET("/", func(r *Ctx) Result {
		sessionWasSet = r.Session != nil
		contextSessionWasSet = GetSession(r.Context()) != nil
		return Text.Result("COOL")
	}, SessionMiddleware(func(_ *Ctx) Result {
		calledCustom = true
		return NoContent
	}))

	meta, err := MockGet(app, "/", r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID)).Discard()
	its.Nil(err)
	its.Equal(http.StatusInternalServerError, meta.StatusCode)
	its.False(sessionWasSet)
	its.False(contextSessionWasSet)
	its.False(calledCustom)
}
