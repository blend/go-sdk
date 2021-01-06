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
	manager.CookieDefaults.SameSite = http.SameSiteLaxMode

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
			if err = opt(am); err != nil {
				return
			}
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

// OptAuthManagerSerializeSessionValueHandler sets a field on an auth manager
func OptAuthManagerSerializeSessionValueHandler(handler AuthManagerSerializeSessionValueHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.SerializeSessionValueHandler = handler
		return nil
	}
}

// OptAuthManagerParseSessionValueHandler sets a field on an auth manager
func OptAuthManagerParseSessionValueHandler(handler AuthManagerParseSessionValueHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.ParseSessionValueHandler = handler
		return nil
	}
}

// OptAuthManagerPersistHandler sets a field on an auth manager
func OptAuthManagerPersistHandler(handler AuthManagerPersistHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.PersistHandler = handler
		return nil
	}
}

// OptAuthManagerFetchHandler sets a field on an auth manager
func OptAuthManagerFetchHandler(handler AuthManagerFetchHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.FetchHandler = handler
		return nil
	}
}

// OptAuthManagerRemoveHandler sets a field on an auth manager
func OptAuthManagerRemoveHandler(handler AuthManagerRemoveHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.RemoveHandler = handler
		return nil
	}
}

// OptAuthManagerValidateHandler sets a field on an auth manager
func OptAuthManagerValidateHandler(handler AuthManagerValidateHandler) AuthManagerOption {
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

// OptAuthManagerPostLoginRedirectHandler sets a field on an auth manager
func OptAuthManagerPostLoginRedirectHandler(handler AuthManagerRedirectHandler) AuthManagerOption {
	return func(am *AuthManager) (err error) {
		am.PostLoginRedirectHandler = handler
		return nil
	}
}

// AuthManagerSerializeSessionValueHandler serializes a session as a string.
type AuthManagerSerializeSessionValueHandler func(context.Context, *Session) (string, error)

// AuthManagerParseSessionValueHandler deserializes a session from a string.
type AuthManagerParseSessionValueHandler func(context.Context, string) (*Session, error)

// AuthManagerPersistHandler saves the session to a stable store.
type AuthManagerPersistHandler func(context.Context, *Session) error

// AuthManagerFetchHandler fetches a session based on a session value.
type AuthManagerFetchHandler func(context.Context, string) (*Session, error)

// AuthManagerRemoveHandler removes a session based on a session value.
type AuthManagerRemoveHandler func(context.Context, string) error

// AuthManagerValidateHandler validates a session.
type AuthManagerValidateHandler func(context.Context, *Session) error

// AuthManagerSessionTimeoutProvider provides a new timeout for a session.
type AuthManagerSessionTimeoutProvider func(*Session) time.Time

// AuthManagerRedirectHandler is a redirect handler.
type AuthManagerRedirectHandler func(*Ctx) *url.URL

// AuthManager is a manager for sessions.
type AuthManager struct {
	CookieDefaults http.Cookie

	SerializeSessionValueHandler AuthManagerSerializeSessionValueHandler
	ParseSessionValueHandler     AuthManagerParseSessionValueHandler

	PersistHandler AuthManagerPersistHandler
	FetchHandler   AuthManagerFetchHandler
	RemoveHandler  AuthManagerRemoveHandler

	ValidateHandler          AuthManagerValidateHandler
	SessionTimeoutProvider   AuthManagerSessionTimeoutProvider
	LoginRedirectHandler     AuthManagerRedirectHandler
	PostLoginRedirectHandler AuthManagerRedirectHandler
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

	// if we're in jwt mode, serialize the jwt.
	if am.SerializeSessionValueHandler != nil {
		sessionValue, err = am.SerializeSessionValueHandler(ctx.Context(), session)
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
	if len(sessionValue) == 0 {
		return nil
	}

	// issue the expiration cookies to the response
	ctx.ExpireCookie(am.CookieDefaults.Name, am.CookieDefaults.Path)
	ctx.Session = nil

	// call the remove handler if one has been provided
	if am.RemoveHandler != nil {
		return am.RemoveHandler(ctx.Context(), sessionValue)
	}
	return nil
}

// VerifySession reads a session value from a request and checks if it's valid.
// It also handles updating a rolling expiry.
//
// It is a pass-through to `VerifyOrExpireSession`
//
// DEPRECATED(1.2021*): this method is deprecated and will be removed.
func (am AuthManager) VerifySession(ctx *Ctx) (*Session, error) {
	return am.VerifyOrExpireSession(ctx)
}

// VerifyOrExpireSession reads a session value from a request and checks if it's valid.
// It also handles updating a rolling expiry.
func (am AuthManager) VerifyOrExpireSession(ctx *Ctx) (session *Session, err error) {
	sessionValue := am.readSessionValue(ctx)
	// validate the sessionValue isn't unset
	if len(sessionValue) == 0 {
		return
	}

	// if we have a separate step to parse the sesion value
	// (i.e. jwt mode) do that now.
	if am.ParseSessionValueHandler != nil {
		session, err = am.ParseSessionValueHandler(ctx.Context(), sessionValue)
		if err != nil {
			if IsErrSessionInvalid(err) {
				_ = am.expire(ctx, sessionValue)
			}
			return
		}
	} else if am.FetchHandler != nil { // if we're in server tracked mode, pull it from whatever backing store we use.
		session, err = am.FetchHandler(ctx.Context(), sessionValue)
		if err != nil {
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
			return nil, err
		}
	}

	if am.SessionTimeoutProvider != nil {
		session.ExpiresUTC = am.SessionTimeoutProvider(session)
		if am.PersistHandler != nil {
			err = am.PersistHandler(ctx.Context(), session)
			if err != nil {
				return nil, err
			}
		}
		am.injectCookie(ctx, sessionValue, session.ExpiresUTC)
	}
	return
}

// VerifySessionValue checks a given session value to see if it's valid.
// It also handles updating a rolling expiry.
func (am AuthManager) VerifySessionValue(ctx *Ctx, sessionValue string) (session *Session, err error) {
	// validate the sessionValue isn't unset
	if len(sessionValue) == 0 {
		return
	}
	// if we have a separate step to parse the sesion value
	// (i.e. jwt mode) do that now.
	if am.ParseSessionValueHandler != nil {
		session, err = am.ParseSessionValueHandler(ctx.Context(), sessionValue)
		if err != nil {
			if IsErrSessionInvalid(err) {
				_ = am.expire(ctx, sessionValue)
			}
			return
		}
	} else if am.FetchHandler != nil { // if we're in server tracked mode, pull it from whatever backing store we use.
		session, err = am.FetchHandler(ctx.Context(), sessionValue)
		if err != nil {
			return
		}
	}

	// if the session is invalid, expire the cookie(s)
	if session == nil || session.IsZero() || session.IsExpired() {
		// return nil whenever the session is invalid
		session = nil
		return
	}

	// call a custom validate handler if one's been provided.
	if am.ValidateHandler != nil {
		err = am.ValidateHandler(ctx.Context(), session)
		if err != nil {
			return nil, err
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

// PostLoginRedirect returns a redirect result for when auth fails and you need to
// send the user to a login page.
func (am AuthManager) PostLoginRedirect(ctx *Ctx) Result {
	if am.PostLoginRedirectHandler != nil {
		redirectTo := am.PostLoginRedirectHandler(ctx)
		if redirectTo != nil {
			return Redirect(redirectTo.String())
		}
	}
	// the default authed redirect is the root.
	return RedirectWithMethod("GET", "/")
}

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

func (am AuthManager) expire(ctx *Ctx, sessionValue string) error {
	ctx.ExpireCookie(am.CookieDefaults.Name, am.CookieDefaults.Path)

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
