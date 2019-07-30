package web

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
)

func TestNewAuthManager(t *testing.T) {
	assert := assert.New(t)

	am, err := NewAuthManager(OptAuthManagerCookieName("X-FOO"))
	assert.Nil(err)
	assert.Equal("X-Foo", am.CookieDefaults.Name)
}

func TestAuthManagerLogin(t *testing.T) {
	assert := assert.New(t)

	am, err := NewLocalAuthManager()
	assert.Nil(err)

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
	serializeHandler := am.SerializeSessionValueHandler
	am.SerializeSessionValueHandler = func(ctx context.Context, session *Session) (string, error) {
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

	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.NotEmpty(session.SessionID)
	assert.NotEmpty(session.RemoteAddr)
	assert.NotEmpty(session.UserAgent)
	assert.Equal("bailey@blend.com", session.UserID)
	assert.True(session.ExpiresUTC.IsZero())
	assert.True(calledPersistHandler)
	assert.True(calledSerializeHandler)
	assert.False(calledRemoveHandler)

	cookies := ReadSetCookies(res.Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.Equal(am.CookieDefaults.Name, cookie.Name)
	assert.Equal(am.CookieDefaults.Path, cookie.Path)
	assert.Equal(session.SessionID, cookie.Value)
}

func TestAuthManagerLogout(t *testing.T) {
	assert := assert.New(t)

	am, err := NewLocalAuthManager()
	assert.Nil(err)

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

	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)

	res = webutil.NewMockResponse(new(bytes.Buffer))
	r = NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))

	assert.Nil(am.Logout(r))
	assert.True(calledRemoveHandler)

	cookies := ReadSetCookies(res.Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.Equal(am.CookieDefaults.Name, cookie.Name)
	assert.Equal(am.CookieDefaults.Path, cookie.Path)
	assert.NotEqual(session.SessionID, cookie.Value, "we should randomize the session cookie on logout")
	assert.True(time.Now().UTC().After(cookie.Expires))
}

func TestAuthManagerVerifySessionParsed(t *testing.T) {
	assert := assert.New(t)

	am, err := NewLocalAuthManager()
	assert.Nil(err)

	var calledParseHandler bool
	am.ParseSessionValueHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledParseHandler = true
		return &Session{UserID: uuid.V4().String(), SessionID: sessionID}, nil
	}

	var calledFetchHandler bool
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		return nil, nil
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, NewSessionID()))
	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.True(session.ExpiresUTC.IsZero())
	assert.True(calledParseHandler)
	assert.False(calledFetchHandler)
	assert.True(calledValidateHandler)
}

func TestAuthManagerVerifySessionFetched(t *testing.T) {
	assert := assert.New(t)

	am, err := NewLocalAuthManager()
	assert.Nil(err)

	var calledFetchHandler bool
	fetchHandler := am.FetchHandler
	am.FetchHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		calledFetchHandler = true
		if fetchHandler == nil {
			return nil, nil
		}
		return fetchHandler(ctx, sessionID)
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
	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.False(calledFetchHandler)
	assert.False(calledValidateHandler)

	r = NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, session.SessionID))
	session, err = am.VerifySession(r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.True(calledFetchHandler)
	assert.True(calledValidateHandler)
}

func TestAuthManagerVerifySessionUnset(t *testing.T) {
	assert := assert.New(t)

	am, err := NewLocalAuthManager()
	assert.Nil(err)

	r := NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))

	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.Nil(session)
}

func TestAuthManagerVerifySessionExpired(t *testing.T) {
	assert := assert.New(t)

	am, err := NewLocalAuthManager()
	assert.Nil(err)

	am.SessionTimeoutProvider = nil
	am.ParseSessionValueHandler = func(ctx context.Context, sessionID string) (*Session, error) {
		return &Session{UserID: uuid.V4().String(), SessionID: sessionID, ExpiresUTC: time.Now().UTC().Add(-time.Hour)}, nil
	}

	var calledValidateHandler bool
	am.ValidateHandler = func(ctx context.Context, session *Session) error {
		calledValidateHandler = true
		return nil
	}

	res := webutil.NewMockResponse(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieDefaults.Name, NewSessionID()))
	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.Nil(session)
	assert.False(calledValidateHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(res.Header())
	assert.NotEmpty(cookies)

	cookie := cookies[0]
	assert.Equal(am.CookieDefaults.Name, cookie.Name)
	assert.Equal(am.CookieDefaults.Path, cookie.Path)
	assert.True(cookie.Expires.Before(time.Now().UTC()), "the cookie should be expired")
}
