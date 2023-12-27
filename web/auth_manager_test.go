/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
)

func Test_MustNewAuthManager(t *testing.T) {
	its := assert.New(t)

	am := MustNewAuthManager(OptAuthManagerCookieName("X-FOO"))
	its.Equal("X-FOO", am.CookieDefaults.Name)

	// test panics
	var recovered interface{}
	func() {
		defer func() {
			r := recover()
			if r != nil {
				recovered = r
			}
		}()
		am = MustNewAuthManager(func(_ *AuthManager) error { return fmt.Errorf("this is just a test") })
	}()
	its.NotNil(recovered)
}

func Test_NewAuthManager(t *testing.T) {
	its := assert.New(t)

	am, err := NewAuthManager()
	its.Nil(err)
	its.Equal(DefaultCookieName, am.CookieDefaults.Name)
	its.Equal(DefaultCookiePath, am.CookieDefaults.Path)
	its.Equal(DefaultCookieHTTPOnly, am.CookieDefaults.HttpOnly)
	its.Equal(DefaultCookieSecure, am.CookieDefaults.Secure)
	its.Equal(DefaultCookieSameSiteMode, am.CookieDefaults.SameSite)

	am, err = NewAuthManager(OptAuthManagerCookieDefaults(http.Cookie{
		Name:		"_FOO_AUTH_",
		Path:		"/admin",
		HttpOnly:	true,
		Secure:		true,
		SameSite:	http.SameSiteLaxMode,
	}))
	its.Nil(err)
	its.Equal("_FOO_AUTH_", am.CookieDefaults.Name)
	its.Equal("/admin", am.CookieDefaults.Path)
	its.Equal(true, am.CookieDefaults.HttpOnly)
	its.Equal(true, am.CookieDefaults.Secure)
	its.Equal(http.SameSiteLaxMode, am.CookieDefaults.SameSite)

	am, err = NewAuthManager(OptAuthManagerCookieName("X-FOO"))
	its.Nil(err)
	its.Equal("X-FOO", am.CookieDefaults.Name)

	am, err = NewAuthManager(OptAuthManagerCookiePath("/foo"))
	its.Nil(err)
	its.Equal("/foo", am.CookieDefaults.Path)

	am, err = NewAuthManager(OptAuthManagerCookieHTTPOnly(true))
	its.Nil(err)
	its.Equal(true, am.CookieDefaults.HttpOnly)

	am, err = NewAuthManager(OptAuthManagerCookieSecure(true))
	its.Nil(err)
	its.Equal(true, am.CookieDefaults.Secure)

	am, err = NewAuthManager(OptAuthManagerCookieSameSite(http.SameSiteLaxMode))
	its.Nil(err)
	its.Equal(http.SameSiteLaxMode, am.CookieDefaults.SameSite)

	am, err = NewAuthManager(OptAuthManagerSerializeHandler(func(context.Context, *Session) (string, error) {
		return "blabla", nil
	}))
	its.Nil(err)
	its.NotNil(am.SerializeHandler)

	am, err = NewAuthManager(OptAuthManagerPersistHandler(func(context.Context, *Session) error {
		return nil
	}))
	its.Nil(err)
	its.NotNil(am.PersistHandler)

	am, err = NewAuthManager(OptAuthManagerFetchHandler(func(context.Context, string) (*Session, error) {
		return &Session{SessionID: "blabla"}, nil
	}))
	its.Nil(err)
	its.NotNil(am.FetchHandler)

	am, err = NewAuthManager(OptAuthManagerRemoveHandler(func(context.Context, string) error {
		return nil
	}))
	its.Nil(err)
	its.NotNil(am.RemoveHandler)

	am, err = NewAuthManager(OptAuthManagerValidateHandler(func(context.Context, *Session) error {
		return nil
	}))
	its.Nil(err)
	its.NotNil(am.ValidateHandler)

	am, err = NewAuthManager(OptAuthManagerSessionTimeoutProvider(func(*Session) time.Time {
		return time.Now().UTC()
	}))
	its.Nil(err)
	its.NotNil(am.SessionTimeoutProvider)

	am, err = NewAuthManager(OptAuthManagerLoginRedirectHandler(func(*Ctx) *url.URL {
		return nil
	}))
	its.Nil(err)
	its.NotNil(am.LoginRedirectHandler)
}

func Test_NewLocalManagerFromCache(t *testing.T) {
	its := assert.New(t)

	lc := NewLocalSessionCache()
	am, err := NewLocalAuthManagerFromCache(lc, OptAuthManagerCookieName("X-FOO"))
	its.Nil(err)
	its.Equal("X-FOO", am.CookieDefaults.Name)

	am, err = NewLocalAuthManagerFromCache(lc, func(_ *AuthManager) error { return fmt.Errorf("this is just a test") })
	its.NotNil(err)
}

func Test_AuthManager_Login(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	sessionExpiresUTC := time.Date(2021, 03, 04, 05, 06, 07, 8, time.UTC)
	var calledSessionTimeoutProvider bool
	am.SessionTimeoutProvider = func(session *Session) time.Time {
		calledSessionTimeoutProvider = true
		return sessionExpiresUTC
	}

	var calledPersistHandler bool
	persistHandler := am.PersistHandler
	am.PersistHandler = func(ctx context.Context, session *Session) error {
		calledPersistHandler = true
		if persistHandler == nil {
			return nil
		}
		return persistHandler(ctx, session)
	}

	var calledSerializeHandler bool
	serializeHandler := am.SerializeHandler
	am.SerializeHandler = func(ctx context.Context, session *Session) (string, error) {
		calledSerializeHandler = true
		if serializeHandler == nil {
			return session.SessionID, nil
		}
		return serializeHandler(ctx, session)
	}

	var calledRemoveHandler bool
	removeHandler := am.RemoveHandler
	am.RemoveHandler = func(ctx context.Context, sessionID string) error {
		calledRemoveHandler = true
		if removeHandler == nil {
			return nil
		}
		return removeHandler(ctx, sessionID)
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.NotEmpty(session.SessionID)
	its.NotEmpty(session.RemoteAddr)
	its.NotEmpty(session.UserAgent)
	its.Equal("example-string@blend.com", session.UserID)
	its.False(session.ExpiresUTC.IsZero())
	its.Equal(sessionExpiresUTC, session.ExpiresUTC)
	its.True(calledPersistHandler)
	its.True(calledSessionTimeoutProvider)
	its.True(calledSerializeHandler)
	its.False(calledRemoveHandler)

	cookies := ReadSetCookies(res.Header())
	its.NotEmpty(cookies)
	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.Equal(session.SessionID, cookie.Value)
}

func Test_AuthManager_Login_persistError(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledPersistHandler bool
	am.PersistHandler = func(ctx context.Context, session *Session) error {
		calledPersistHandler = true
		return fmt.Errorf("this is just a test")
	}
	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

	session, err := am.Login("example-string@blend.com", r)
	its.NotNil(err)
	its.True(calledPersistHandler)
	its.Equal("this is just a test", err.Error())
	its.Nil(session)

	cookies := ReadSetCookies(res.Header())
	its.Empty(cookies)
}

func Test_AuthManager_Login_serializeError(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledPersistHandler bool
	persistHandler := am.PersistHandler
	am.PersistHandler = func(ctx context.Context, session *Session) error {
		calledPersistHandler = true
		if persistHandler == nil {
			return nil
		}
		return persistHandler(ctx, session)
	}

	var calledSerializeHandler bool
	am.SerializeHandler = func(ctx context.Context, session *Session) (string, error) {
		calledSerializeHandler = true
		return "", fmt.Errorf("this is a serialize error")
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

	session, err := am.Login("example-string@blend.com", r)
	its.NotNil(err)
	its.True(calledPersistHandler)
	its.True(calledSerializeHandler)
	its.Equal("this is a serialize error", err.Error())
	its.Nil(session)

	cookies := ReadSetCookies(res.Header())
	its.Empty(cookies)
}

func Test_AuthManager_Logout(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledRemoveHandler bool
	removeHandler := am.RemoveHandler
	am.RemoveHandler = func(ctx context.Context, sessionID string) error {
		calledRemoveHandler = true
		if removeHandler == nil {
			return nil
		}
		return removeHandler(ctx, sessionID)
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)

	res = webutil.NewMockResponse(new(bytes.Buffer))
	r = NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))

	its.Nil(am.Logout(r))
	its.True(calledRemoveHandler)

	cookies := ReadSetCookies(res.Header())
	its.NotEmpty(cookies)
	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.NotEqual(session.SessionID, cookie.Value, "we should randomize the session cookie on logout")
	its.True(time.Now().UTC().After(cookie.Expires))
}

func Test_AuthManager_Logout_sessionValueUnset(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))
	its.Nil(am.Logout(r))

	cookies := ReadSetCookies(res.Header())
	its.Empty(cookies)
}

func Test_AuthManager_Logout_removeHandlerUnset(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	am.RemoveHandler = nil
	its.Nil(err)

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)

	res = webutil.NewMockResponse(new(bytes.Buffer))
	r = NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))

	its.Nil(am.Logout(r))

	cookies := ReadSetCookies(res.Header())
	its.NotEmpty(cookies)
}

func Test_AuthManager_VerifyOrExtendSession(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledRestoreHandler bool
	restoreHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledRestoreHandler = true
		if restoreHandler == nil {
			return nil, nil
		}
		return restoreHandler(ctx, sessionID)
	}

	var calledValidateHandler bool
	validateHandler := am.ValidateHandler
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		if validateHandler == nil {
			return nil
		}
		return validateHandler(ctx, session)
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.False(calledRestoreHandler)
	its.False(calledValidateHandler)

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.NotNil(session)
	its.True(calledRestoreHandler)
	its.True(calledValidateHandler)
}

func Test_AuthManager_VerifyOrExtendSession_sessionUnset(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledRestoreHandler bool
	restoreHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledRestoreHandler = true
		if restoreHandler == nil {
			return nil, nil
		}
		return restoreHandler(ctx, sessionID)
	}

	var calledValidateHandler bool
	validateHandler := am.ValidateHandler
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		if validateHandler == nil {
			return nil
		}
		return validateHandler(ctx, session)
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.Nil(session)
	its.False(calledRestoreHandler)
	its.False(calledValidateHandler)
}

func Test_AuthManager_VerifyOrExtendSession_fetchHandlerUnset(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))

	session, err := am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.Nil(session)
}

func Test_AuthManager_VerifyOrExtendSession_fetchErrSessionInvalid(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledFetchHandler bool
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		return nil, ErrSessionIDEmpty
	}

	var calledValidateHandler bool
	validateHandler := am.ValidateHandler
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return validateHandler(ctx, session)
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.False(calledFetchHandler)
	its.False(calledValidateHandler)

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifyOrExtendSession(r)
	its.NotNil(err)
	its.Equal(ErrSessionIDEmpty, err)
	its.Nil(session)
	its.True(calledFetchHandler)
	its.False(calledValidateHandler)

	cookies := ReadSetCookies(r.Response.Header())
	its.NotEmpty(cookies)
	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.True(time.Now().UTC().After(cookie.Expires))
}

func Test_AuthManager_VerifyOrExtendSession_sessionExpired(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	am.SessionTimeoutProvider = nil
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		return &Session{UserID: uuid.V4().String(), SessionID: sessionID, ExpiresUTC: time.Now().UTC().Add(-time.Hour)}, nil
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, NewSessionID()))
	session, err := am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.Nil(session)
	its.False(calledValidateHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(res.Header())
	its.NotEmpty(cookies)

	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.True(cookie.Expires.Before(time.Now().UTC()), "the cookie should be expired")
}

func Test_AuthManager_VerifyOrExtendSession_sessionExpired_nil(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	am.SessionTimeoutProvider = nil
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		return nil, nil
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, NewSessionID()))
	session, err := am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.Nil(session)
	its.False(calledValidateHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(res.Header())
	its.NotEmpty(cookies)

	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.True(cookie.Expires.Before(time.Now().UTC()), "the cookie should be expired")
}

func Test_AuthManager_VerifyOrExtendSession_sessionExpired_zero(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	am.SessionTimeoutProvider = nil
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		return &Session{}, nil
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, NewSessionID()))
	session, err := am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.Nil(session)
	its.False(calledValidateHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(res.Header())
	its.NotEmpty(cookies)

	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.True(cookie.Expires.Before(time.Now().UTC()), "the cookie should be expired")
}

func Test_AuthManager_VerifyOrExtendSession_failsValidation(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledFetchHandler bool
	fetchHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		return fetchHandler(ctx, sessionID)
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return fmt.Errorf("this is just a test")
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.False(calledFetchHandler)
	its.False(calledValidateHandler)

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifyOrExtendSession(r)
	its.NotNil(err)
	its.Nil(session)
	its.True(calledFetchHandler)
	its.True(calledValidateHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(r.Response.Header())
	// for now, we should not expire the cookie on a validatin failure
	its.Empty(cookies)
}

func Test_AuthManager_VerifyOrExtendSession_sessionTimeout_unchanged(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledPersistHandler bool
	persistHandler := am.PersistHandler
	am.PersistHandler = func(ctx context.Context, session *Session) error {
		calledPersistHandler = true
		return persistHandler(ctx, session)
	}

	var calledFetchHandler bool
	fetchHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		return fetchHandler(ctx, sessionID)
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}
	var calledSessionTimeoutProvider bool
	am.SessionTimeoutProvider = func(session *Session) time.Time {
		calledSessionTimeoutProvider = true
		return session.ExpiresUTC
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.False(calledFetchHandler)
	its.False(calledValidateHandler)
	its.True(calledPersistHandler)
	its.True(calledSessionTimeoutProvider)

	calledPersistHandler = false
	calledSessionTimeoutProvider = false

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.NotNil(session)
	its.True(calledFetchHandler)
	its.True(calledValidateHandler)
	its.False(calledPersistHandler)
	its.True(calledSessionTimeoutProvider)
}

func Test_AuthManager_VerifyOrExtendSession_sessionTimeout_changed(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledFetchHandler bool
	fetchHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		return fetchHandler(ctx, sessionID)
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.False(calledFetchHandler)
	its.False(calledValidateHandler)

	var calledPersistHandler bool
	persistHandler := am.PersistHandler
	am.PersistHandler = func(ctx context.Context, session *Session) error {
		calledPersistHandler = true
		return persistHandler(ctx, session)
	}

	expiresUTC := time.Now().UTC()
	var calledSessionTimeoutProvider bool
	am.SessionTimeoutProvider = func(session *Session) time.Time {
		calledSessionTimeoutProvider = true
		return expiresUTC
	}

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifyOrExtendSession(r)
	its.Nil(err)
	its.NotNil(session)
	its.True(calledFetchHandler)
	its.True(calledValidateHandler)
	its.True(calledPersistHandler)
	its.True(calledSessionTimeoutProvider)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(r.Response.Header())
	its.NotEmpty(cookies)

	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.InTimeDelta(expiresUTC, cookie.Expires, time.Second, "the cookie have an expiration")
}

func Test_AuthManager_VerifyOrExpireSession_sessionTimeout_persistError(t *testing.T) {
	its := assert.New(t)

	am, err := NewLocalAuthManager()
	its.Nil(err)

	var calledFetchHandler bool
	fetchHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		return fetchHandler(ctx, sessionID)
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("example-string@blend.com", r)
	its.Nil(err)
	its.NotNil(session)
	its.False(calledFetchHandler)
	its.False(calledValidateHandler)

	var calledPersistHandler bool
	am.PersistHandler = func(ctx context.Context, session *Session) error {
		calledPersistHandler = true
		return fmt.Errorf("this is just a test")
	}

	expiresUTC := time.Now().UTC()
	var calledSessionTimeoutProvider bool
	am.SessionTimeoutProvider = func(session *Session) time.Time {
		calledSessionTimeoutProvider = true
		return expiresUTC
	}

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifyOrExtendSession(r)
	its.NotNil(err)
	its.Equal("this is just a test", err.Error())
	its.Nil(session)
	its.True(calledFetchHandler)
	its.True(calledValidateHandler)
	its.True(calledSessionTimeoutProvider)
	its.True(calledPersistHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(r.Response.Header())
	its.Empty(cookies)
}

func Test_AuthManager_LoginRedirect_loginRedirectHandlerUnset(t *testing.T) {
	its := assert.New(t)

	am := AuthManager{}
	ctx := MockCtx(http.MethodGet, "/api/foo/bar", OptCtxDefaultProvider(Text))
	res := am.LoginRedirect(ctx)
	its.NotNil(res)
	typed, ok := res.(*RawResult)
	its.True(ok)
	its.Equal(http.StatusUnauthorized, typed.StatusCode)
}

func Test_AuthManager_LoginRedirect_loginRedirectHandler(t *testing.T) {
	its := assert.New(t)

	am := AuthManager{
		LoginRedirectHandler: func(r *Ctx) *url.URL {
			return &url.URL{
				Path: "/foo",
			}
		},
	}
	ctx := MockCtx(http.MethodGet, "/api/foo/bar", OptCtxDefaultProvider(Text))
	res := am.LoginRedirect(ctx)
	its.NotNil(res)
	typed, ok := res.(*RedirectResult)
	its.True(ok)
	its.Empty(typed.Method)
	its.Equal("/foo", typed.RedirectURI)
}

func Test_AuthManager_expire(t *testing.T) {
	its := assert.New(t)

	expectedSessionID := uuid.V4().String()
	cookieName := "my-auth-cookie"

	var didCallRemoveHandler, sessionIDCorrect bool
	am := AuthManager{
		CookieDefaults: http.Cookie{
			Name: cookieName,
		},
		RemoveHandler: func(c context.Context, sessionID string) error {
			didCallRemoveHandler = true
			sessionIDCorrect = sessionID == expectedSessionID
			return nil
		},
	}
	ctx := MockCtx(http.MethodGet, "/api/foo/bar", OptCtxDefaultProvider(Text), OptCtxCookieValue(cookieName, expectedSessionID))
	err := am.expire(ctx, expectedSessionID)
	its.Nil(err)
	its.True(didCallRemoveHandler)
	its.True(sessionIDCorrect)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(ctx.Response.Header())
	its.NotEmpty(cookies)

	cookie := cookies[0]
	its.Equal(am.CookieDefaults.Name, cookie.Name)
	its.Equal(am.CookieDefaults.Path, cookie.Path)
	its.True(cookie.Expires.Before(time.Now().UTC()), "the cookie should be expired")
}

func Test_AuthManager_expire_removeError(t *testing.T) {
	its := assert.New(t)

	var didCallRemoveHandler bool
	am := AuthManager{
		RemoveHandler: func(c context.Context, sessionID string) error {
			didCallRemoveHandler = true
			return fmt.Errorf("this is just a test")
		},
	}
	ctx := MockCtx(http.MethodGet, "/api/foo/bar", OptCtxDefaultProvider(Text))
	err := am.expire(ctx, uuid.V4().String())
	its.NotNil(err)
	its.Equal("this is just a test", err.Error())
	its.True(didCallRemoveHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(ctx.Response.Header())
	its.Empty(cookies)
}
