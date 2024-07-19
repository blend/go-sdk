/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

//github:codeowner @blend/infosec

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// MustNewAuthManager returns a new auth manager with a given set of options but panics on error.
func MustNewAuthManager(options ...AuthManagerOption) AuthManager {
	am, err := NewAuthManager(options...)
	if err != nil {
		panic(err)
	}
	return am
}

// NewAuthManager returns a new auth manager from a given config.
// For remote mode, you must provide a fetch, persist, and remove handler, and optionally a login redirect handler.
func NewAuthManager(options ...AuthManagerOption) (manager AuthManager, err error) {
	manager.CookieDefaults.Name = DefaultCookieName
	manager.CookieDefaults.Path = DefaultCookiePath
	manager.CookieDefaults.Secure = DefaultCookieSecure
	manager.CookieDefaults.HttpOnly = DefaultCookieHTTPOnly
	manager.CookieDefaults.SameSite = DefaultCookieSameSiteMode

	for _, opt := range options {
		if err = opt(&manager); err != nil {
			return
		}
	}
	return
}

// NewLocalAuthManager returns a new locally cached session manager.
// It saves sessions to a local store.
func NewLocalAuthManager(options ...AuthManagerOption) (AuthManager, error) {
	return NewLocalAuthManagerFromCache(NewLocalSessionCache(), options...)
}

// NewLocalAuthManagerFromCache returns a new locally cached session manager that saves sessions to the cache provided
func NewLocalAuthManagerFromCache(cache *LocalSessionCache, options ...AuthManagerOption) (manager AuthManager, err error) {
	manager, err = NewAuthManager(options...)
	if err != nil {
		return
	}
	manager.PersistHandler = cache.PersistHandler
	manager.FetchHandler = cache.FetchHandler
	manager.RemoveHandler = cache.RemoveHandler
	return
}

// AuthManagerOption is a variadic option for auth managers.
type AuthManagerOption func(*AuthManager) error

// OptAuthManagerFromConfig returns an auth manager from a config.
func OptAuthManagerFromConfig(cfg Config) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		opts := []AuthManagerOption{
			OptAuthManagerCookieSecure(cfg.CookieSecureOrDefault()),
			OptAuthManagerCookieHTTPOnly(cfg.CookieHTTPOnlyOrDefault()),
			OptAuthManagerCookieName(cfg.CookieNameOrDefault()),
			OptAuthManagerCookiePath(cfg.CookiePathOrDefault()),
			OptAuthManagerCookieDomain(cfg.CookieDomainOrDefault()),
			OptAuthManagerCookieSameSite(cfg.CookieSameSiteOrDefault()),
			OptAuthManagerSessionTimeoutProvider(SessionTimeoutProvider(!cfg.SessionTimeoutIsRelative, cfg.SessionTimeoutOrDefault())),
		}
		for _, opt := range opts {
			// NOTE(wc): none of the options above produce an error
			// it is safe to ignore the error produced
			// by the option call.
			_ = opt(am)
		}
		return
	}
}

// OptAuthManagerCookieDefaults sets a field on an auth manager
func OptAuthManagerCookieDefaults(cookie http.Cookie) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults = cookie
		return nil
	}
}

// OptAuthManagerCookieSecure sets a field on an auth manager
func OptAuthManagerCookieSecure(secure bool) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults.Secure = secure
		return nil
	}
}

// OptAuthManagerCookieHTTPOnly sets a field on an auth manager
func OptAuthManagerCookieHTTPOnly(httpOnly bool) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults.HttpOnly = httpOnly
		return nil
	}
}

// OptAuthManagerCookieName sets a field on an auth manager
func OptAuthManagerCookieName(cookieName string) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults.Name = cookieName
		return nil
	}
}

// OptAuthManagerCookiePath sets a field on an auth manager
func OptAuthManagerCookiePath(cookiePath string) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults.Path = cookiePath
		return nil
	}
}

// OptAuthManagerCookieDomain sets a field on an auth manager
func OptAuthManagerCookieDomain(domain string) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults.Domain = domain
		return nil
	}
}

// OptAuthManagerCookieSameSite sets a field on an auth manager
func OptAuthManagerCookieSameSite(sameSite http.SameSite) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.CookieDefaults.SameSite = sameSite
		return nil
	}
}

// OptAuthManagerSerializeHandler sets a field on an auth manager
func OptAuthManagerSerializeHandler(handler AuthManagerSerializeSessionHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.SerializeHandler = handler
		return nil
	}
}

// OptAuthManagerPersistHandler sets a field on an auth manager
func OptAuthManagerPersistHandler(handler AuthManagerPersistSessionHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.PersistHandler = handler
		return nil
	}
}

// OptAuthManagerFetchHandler sets a field on an auth manager
func OptAuthManagerFetchHandler(handler AuthManagerFetchSessionHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.FetchHandler = handler
		return nil
	}
}

// OptAuthManagerRemoveHandler sets a field on an auth manager
func OptAuthManagerRemoveHandler(handler AuthManagerRemoveSessionHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.RemoveHandler = handler
		return nil
	}
}

// OptAuthManagerValidateHandler sets a field on an auth manager
func OptAuthManagerValidateHandler(handler AuthManagerValidateSessionHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.ValidateHandler = handler
		return nil
	}
}

// OptAuthManagerSessionTimeoutProvider sets a field on an auth manager
func OptAuthManagerSessionTimeoutProvider(handler AuthManagerSessionTimeoutProvider) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.SessionTimeoutProvider = handler
		return nil
	}
}

// OptAuthManagerLoginRedirectHandler sets a field on an auth manager
func OptAuthManagerLoginRedirectHandler(handler AuthManagerRedirectHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.LoginRedirectHandler = handler
		return nil
	}
}

// AuthManagerSerializeSessionHandler serializes a session as a string.
type AuthManagerSerializeSessionHandler func(context.Context, *Session) (string, error)

// AuthManagerPersistSessionHandler saves the session to a stable store.
type AuthManagerPersistSessionHandler func(context.Context, *Session) error

// AuthManagerFetchSessionHandler restores a session based on a session value.
type AuthManagerFetchSessionHandler func(context.Context, string) (*Session, error)

// AuthManagerRemoveSessionHandler removes a session based on a session value.
type AuthManagerRemoveSessionHandler func(context.Context, string) error

// AuthManagerValidateSessionHandler validates a session.
type AuthManagerValidateSessionHandler func(context.Context, *Session) error

// AuthManagerSessionTimeoutProvider provides a new timeout for a session.
type AuthManagerSessionTimeoutProvider func(*Session) time.Time

// AuthManagerRedirectHandler is a redirect handler.
type AuthManagerRedirectHandler func(*Ctx) *url.URL

// AuthManager is a manager for sessions.
type AuthManager struct {
	CookieDefaults http.Cookie

	// PersistHandler is called to both create and to update a session in a persistent store.
	PersistHandler AuthManagerPersistSessionHandler
	// SerializeSessionHandler if set, is called to serialize the session
	// as a session cookie value.
	SerializeHandler AuthManagerSerializeSessionHandler
	// FetchSessionHandler is called if set to restore a session from a string session identifier.
	FetchHandler AuthManagerFetchSessionHandler
	// Remove handler is called on logout to remove a session from a persistent store.
	// It is called during `Logout` to remove logged out sessions.
	RemoveHandler AuthManagerRemoveSessionHandler
	// ValidateHandler is called after a session is retored to make sure it's still valid.
	ValidateHandler AuthManagerValidateSessionHandler
	// SessionTimeoutProvider is called to create a variable session expiry.
	SessionTimeoutProvider AuthManagerSessionTimeoutProvider

	// LoginRedirectHandler redirects an unauthenticated user to the login page.
	LoginRedirectHandler AuthManagerRedirectHandler
}

// --------------------------------------------------------------------------------
// Methods
// --------------------------------------------------------------------------------

// Login logs a userID in.
func (am AuthManager) Login(userID string, ctx *Ctx) (session *Session, err error) {
	// create a new session value
	sessionValue := NewSessionID()
	// userID and sessionID are required
	session = NewSession(userID, sessionValue)
	if am.SessionTimeoutProvider != nil {
		session.ExpiresUTC = am.SessionTimeoutProvider(session)
	}
	session.UserAgent = webutil.GetUserAgent(ctx.Request)
	session.RemoteAddr = webutil.GetRemoteAddr(ctx.Request)

	// call the perist handler if one's been provided
	if am.PersistHandler != nil {
		err = am.PersistHandler(ctx.Context(), session)
		if err != nil {
			return nil, err
		}
	}

	// call the serialize handler if one's been provided
	if am.SerializeHandler != nil {
		sessionValue, err = am.SerializeHandler(ctx.Context(), session)
		if err != nil {
			return nil, err
		}
	}

	// inject cookies into the response
	am.injectCookie(ctx, sessionValue, session.ExpiresUTC)
	return session, nil
}

// Logout unauthenticates a session.
func (am AuthManager) Logout(ctx *Ctx) error {
	sessionValue := am.readSessionValue(ctx)
	// validate the sessionValue isn't unset
	if sessionValue == "" {
		return nil
	}
	// zero out the context session as a precaution
	ctx.Session = nil
	// issue the expiration cookies to the response
	// and call the remove handler
	return am.expire(ctx, sessionValue)
}

// VerifySession pulls the session cookie off the request, and validates
// it represents a valid session.
func (am AuthManager) VerifySession(ctx *Ctx) (sessionValue string, session *Session, err error) {
	sessionValue = am.readSessionValue(ctx)
	// validate the sessionValue is set
	if len(sessionValue) == 0 {
		return
	}

	// if we have a restore handler, call it.
	if am.FetchHandler != nil {
		session, err = am.FetchHandler(ctx.Context(), sessionValue)
		if err != nil {
			if IsErrSessionInvalid(err) {
				_ = am.expire(ctx, sessionValue)
			}
			return
		}
	}

	// if the session is invalid, expire the cookie(s)
	if session == nil || session.IsZero() || session.IsExpired() {
		// return nil whenever the session is invalid
		session = nil
		err = am.expire(ctx, sessionValue)
		return
	}

	// call a custom validate handler if one's been provided.
	if am.ValidateHandler != nil {
		err = am.ValidateHandler(ctx.Context(), session)
		if err != nil {
			session = nil
			return
		}
	}
	return
}

// VerifyOrExtendSession reads a session value from a request and checks if it's valid.
// It also handles updating a rolling expiry.
func (am AuthManager) VerifyOrExtendSession(ctx *Ctx) (session *Session, err error) {
	var sessionValue string
	sessionValue, session, err = am.VerifySession(ctx)
	if session == nil || err != nil {
		return
	}

	if am.SessionTimeoutProvider != nil {
		existingExpiresUTC := session.ExpiresUTC
		session.ExpiresUTC = am.SessionTimeoutProvider(session)

		// if session expiry has changed
		if existingExpiresUTC != session.ExpiresUTC {
			// if we have a persist handler
			// call it to reflect the updated session timeout.
			if am.PersistHandler != nil {
				err = am.PersistHandler(ctx.Context(), session)
				if err != nil {
					return nil, err
				}
			}

			// inject the (updated) cookie
			am.injectCookie(ctx, sessionValue, session.ExpiresUTC)
		}
	}
	return
}

// LoginRedirect returns a redirect result for when auth fails and you need to
// send the user to a login page.
func (am AuthManager) LoginRedirect(ctx *Ctx) Result {
	if am.LoginRedirectHandler != nil {
		redirectTo := am.LoginRedirectHandler(ctx)
		if redirectTo != nil {
			return Redirect(redirectTo.String())
		}
	}
	return ctx.DefaultProvider.NotAuthorized()
}

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

func (am AuthManager) expire(ctx *Ctx, sessionValue string) error {
	// issue the cookie expiration.
	am.expireCookie(ctx)

	// if we have a remove handler and the sessionID is set
	if am.RemoveHandler != nil {
		err := am.RemoveHandler(ctx.Context(), sessionValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// InjectCookie injects a session cookie into the context.
func (am AuthManager) injectCookie(ctx *Ctx, value string, expire time.Time) {
	http.SetCookie(ctx.Response, &http.Cookie{
		Value:    value,
		Expires:  expire,
		Name:     am.CookieDefaults.Name,
		Path:     am.CookieDefaults.Path,
		Domain:   am.CookieDefaults.Domain,
		HttpOnly: am.CookieDefaults.HttpOnly,
		Secure:   am.CookieDefaults.Secure,
		SameSite: am.CookieDefaults.SameSite,
	})
}

// expireCookie expires the session cookie.
func (am AuthManager) expireCookie(ctx *Ctx) {
	http.SetCookie(ctx.Response, &http.Cookie{
		Value: NewSessionID(),
		// MaxAge<0 means delete cookie now, and is equivalent to
		// the literal cookie header content 'Max-Age: 0'
		MaxAge:   -1,
		Name:     am.CookieDefaults.Name,
		Path:     am.CookieDefaults.Path,
		Domain:   am.CookieDefaults.Domain,
		HttpOnly: am.CookieDefaults.HttpOnly,
		Secure:   am.CookieDefaults.Secure,
		SameSite: am.CookieDefaults.SameSite,
	})

}

// cookieValue reads a param from a given request context from either the cookies or headers.
func (am AuthManager) cookieValue(name string, ctx *Ctx) (output string) {
	if cookie := ctx.Cookie(name); cookie != nil {
		output = cookie.Value
	}
	return
}

// ReadSessionID reads a session id from a given request context.
func (am AuthManager) readSessionValue(ctx *Ctx) string {
	return am.cookieValue(am.CookieDefaults.Name, ctx)
}
