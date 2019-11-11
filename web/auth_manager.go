package web

import (
	"context"
	"net/url"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// AuthManagerMode is an auth manager mode.
type AuthManagerMode string

const (
	// AuthManagerModeJWT is the jwt auth mode.
	AuthManagerModeJWT AuthManagerMode = "jwt"
	// AuthManagerModeServer is the server managed auth mode.
	AuthManagerModeServer AuthManagerMode = "server"
	// AuthManagerModeLocal is the local cache auth mode.
	AuthManagerModeLocal AuthManagerMode = "cached"
)

// NewAuthManagerFromConfig returns a new auth manager from a given config.
func NewAuthManagerFromConfig(cfg *Config) (manager *AuthManager) {
	switch cfg.GetAuthManagerMode() {
	case AuthManagerModeJWT:
		manager = NewJWTAuthManager(cfg.GetAuthSecret())
	case AuthManagerModeLocal: // local should only be used for debugging.
		manager = NewLocalAuthManager()
	case AuthManagerModeServer:
		manager = NewServerAuthManager()
	default:
		panic("invalid auth manager mode")
	}

	return manager.WithCookieHTTPSOnly(cfg.GetCookieHTTPSOnly()).
		WithCookieName(cfg.GetCookieName()).
		WithCookiePath(cfg.GetCookiePath()).
		WithSessionTimeoutProvider(SessionTimeoutProvider(cfg.GetSessionTimeoutIsAbsolute(), cfg.GetSessionTimeout()))
}

// NewLocalAuthManager returns a new locally cached session manager that saves sessions to the cache provided
func NewLocalAuthManagerFromCache(cache *LocalSessionCache) *AuthManager {
	return &AuthManager{
		persistHandler: cache.PersistHandler,
		fetchHandler:   cache.FetchHandler,
		removeHandler:  cache.RemoveHandler,
		cookieName:     DefaultCookieName,
		cookiePath:     DefaultCookiePath,
	}
}

// NewLocalAuthManager returns a new locally cached session manager.
// It saves sessions to a local store.
func NewLocalAuthManager() *AuthManager {
	cache := NewLocalSessionCache()
	return &AuthManager{
		persistHandler: cache.PersistHandler,
		fetchHandler:   cache.FetchHandler,
		removeHandler:  cache.RemoveHandler,
		cookieName:     DefaultCookieName,
		cookiePath:     DefaultCookiePath,
	}
}

// NewJWTAuthManager returns a new jwt session manager.
// It issues JWT tokens to identify users.
func NewJWTAuthManager(key []byte) *AuthManager {
	jwtm := NewJWTManager(key)
	return &AuthManager{
		serializeSessionValueHandler: jwtm.SerializeSessionValueHandler,
		parseSessionValueHandler:     jwtm.ParseSessionValueHandler,
		cookieName:                   DefaultCookieName,
		cookiePath:                   DefaultCookiePath,
		sessionTimeoutProvider:       SessionTimeoutProviderAbsolute(DefaultSessionTimeout),
	}
}

// NewServerAuthManager returns a new server auth manager.
// You should set the `FetchHandler`, the `PersistHandler` and the `RemoveHandler`.
func NewServerAuthManager() *AuthManager {
	return &AuthManager{
		cookieName: DefaultCookieName,
		cookiePath: DefaultCookiePath,
	}
}

// AuthManagerSerializeSessionValueHandler serializes a session as a string.
type AuthManagerSerializeSessionValueHandler func(context.Context, *Session, State) (string, error)

// AuthManagerParseSessionValueHandler deserializes a session from a string.
type AuthManagerParseSessionValueHandler func(context.Context, string, State) (*Session, error)

// AuthManagerPersistHandler saves the session to a stable store.
type AuthManagerPersistHandler func(context.Context, *Session, State) error

// AuthManagerFetchHandler fetches a session based on a session value.
type AuthManagerFetchHandler func(context.Context, string, State) (*Session, error)

// AuthManagerRemoveHandler removes a session based on a session value.
type AuthManagerRemoveHandler func(context.Context, string, State) error

// AuthManagerValidateHandler validates a session.
type AuthManagerValidateHandler func(context.Context, *Session, State) error

// AuthManagerSessionTimeoutProvider provides a new timeout for a session.
type AuthManagerSessionTimeoutProvider func(*Session) time.Time

// AuthManagerRedirectHandler is a redirect handler.
type AuthManagerRedirectHandler func(*Ctx) *url.URL

// AuthManager is a manager for sessions.
type AuthManager struct {
	serializeSessionValueHandler AuthManagerSerializeSessionValueHandler
	parseSessionValueHandler     AuthManagerParseSessionValueHandler

	// these generally apply to server or local modes.
	persistHandler AuthManagerPersistHandler
	fetchHandler   AuthManagerFetchHandler
	removeHandler  AuthManagerRemoveHandler

	// these generally apply to any mode.
	validateHandler          AuthManagerValidateHandler
	sessionTimeoutProvider   AuthManagerSessionTimeoutProvider
	loginRedirectHandler     AuthManagerRedirectHandler
	postLoginRedirectHandler AuthManagerRedirectHandler

	cookieName      string
	cookiePath      string
	cookieHTTPSOnly bool
}

// --------------------------------------------------------------------------------
// Methods
// --------------------------------------------------------------------------------

// Login logs a userID in.
func (am *AuthManager) Login(userID string, ctx *Ctx) (session *Session, err error) {
	// create a new session value
	sessionValue := NewSessionID()
	// userID and sessionID are required
	session = NewSession(userID, sessionValue)
	if am.sessionTimeoutProvider != nil {
		session.ExpiresUTC = am.sessionTimeoutProvider(session)
	}
	session.UserAgent = webutil.GetUserAgent(ctx.request)
	session.RemoteAddr = webutil.GetRemoteAddr(ctx.request)

	// call the perist handler if one's been provided
	if am.persistHandler != nil {
		err = am.persistHandler(ctx.Context(), session, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	// if we're in jwt mode, serialize the jwt.
	if am.serializeSessionValueHandler != nil {
		sessionValue, err = am.serializeSessionValueHandler(ctx.Context(), session, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	// inject cookies into the response
	am.injectCookie(ctx, am.CookieName(), sessionValue, session.ExpiresUTC)
	return session, nil
}

// Logout unauthenticates a session.
func (am *AuthManager) Logout(ctx *Ctx) error {
	sessionValue := am.readSessionValue(ctx)
	// validate the sessionValue isn't unset
	if len(sessionValue) == 0 {
		return nil
	}

	// issue the expiration cookies to the response
	ctx.ExpireCookie(am.CookieName(), am.CookiePath())
	// nil out the current session in the ctx
	ctx.WithSession(nil)

	// call the remove handler if one has been provided
	if am.removeHandler != nil {
		return am.removeHandler(ctx.Context(), sessionValue, ctx.state)
	}
	return nil
}

// VerifySession checks a sessionID to see if it's valid.
// It also handles updating a rolling expiry.
func (am *AuthManager) VerifySession(ctx *Ctx) (session *Session, err error) {
	// pull the sessionID off the request
	sessionValue := am.readSessionValue(ctx)
	// validate the sessionValue isn't unset
	if len(sessionValue) == 0 {
		return
	}

	// if we have a separate step to parse the sesion value
	// (i.e. jwt mode) do that now.
	if am.parseSessionValueHandler != nil {
		session, err = am.parseSessionValueHandler(ctx.Context(), sessionValue, ctx.state)
		if err != nil {
			if IsErrSessionInvalid(err) {
				am.expire(ctx, sessionValue)
			}
			return
		}
	} else if am.fetchHandler != nil { // if we're in server tracked mode, pull it from whatever backing store we use.
		session, err = am.fetchHandler(ctx.Context(), sessionValue, ctx.state)
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
	if am.validateHandler != nil {
		err = am.validateHandler(ctx.Context(), session, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	if am.sessionTimeoutProvider != nil {
		session.ExpiresUTC = am.sessionTimeoutProvider(session)
		if am.persistHandler != nil {
			err = am.persistHandler(ctx.Context(), session, ctx.state)
			if err != nil {
				return nil, err
			}
		}
		am.injectCookie(ctx, am.CookieName(), sessionValue, session.ExpiresUTC)
	}
	return
}

// LoginRedirect returns a redirect result for when auth fails and you need to
// send the user to a login page.
func (am *AuthManager) LoginRedirect(ctx *Ctx) Result {
	if am.loginRedirectHandler != nil {
		redirectTo := am.loginRedirectHandler(ctx)
		if redirectTo != nil {
			return ctx.Redirect(redirectTo.String())
		}
	}
	return ctx.DefaultResultProvider().NotAuthorized()
}

// PostLoginRedirect returns a redirect result for when auth fails and you need to
// send the user to a login page.
func (am *AuthManager) PostLoginRedirect(ctx *Ctx) Result {
	if am.postLoginRedirectHandler != nil {
		redirectTo := am.postLoginRedirectHandler(ctx)
		if redirectTo != nil {
			return ctx.Redirect(redirectTo.String())
		}
	}
	// the default authed redirect is the root.
	return ctx.RedirectWithMethod("GET", "/")
}

// --------------------------------------------------------------------------------
// Properties
// --------------------------------------------------------------------------------

// WithSessionTimeoutProvider sets the session timeout provider.
func (am *AuthManager) WithSessionTimeoutProvider(timeoutProvider func(*Session) time.Time) *AuthManager {
	am.sessionTimeoutProvider = timeoutProvider
	return am
}

// SessionTimeoutProvider returns the session timeout provider.
func (am *AuthManager) SessionTimeoutProvider() func(*Session) time.Time {
	return am.sessionTimeoutProvider
}

// WithCookieHTTPSOnly sets if we should issue cookies with the HTTPS flag on.
func (am *AuthManager) WithCookieHTTPSOnly(isHTTPSOnly bool) *AuthManager {
	am.cookieHTTPSOnly = isHTTPSOnly
	return am
}

// CookiesHTTPSOnly returns if the cookie is for only https connections.
func (am *AuthManager) CookiesHTTPSOnly() bool {
	return am.cookieHTTPSOnly
}

// WithCookieName sets the cookie name.
func (am *AuthManager) WithCookieName(paramName string) *AuthManager {
	am.cookieName = paramName
	return am
}

// CookieName returns the session param name.
func (am *AuthManager) CookieName() string {
	return am.cookieName
}

// WithCookiePath sets the cookie path.
func (am *AuthManager) WithCookiePath(path string) *AuthManager {
	am.cookiePath = path
	return am
}

// CookiePath returns the session param path.
func (am *AuthManager) CookiePath() string {
	if len(am.cookiePath) == 0 {
		return DefaultCookiePath
	}
	return am.cookiePath
}

// WithSerializeSessionValueHandler sets the serialize session value handler.
func (am *AuthManager) WithSerializeSessionValueHandler(handler AuthManagerSerializeSessionValueHandler) *AuthManager {
	am.serializeSessionValueHandler = handler
	return am
}

// SerializeSessionValueHandler returns the serialize session value handler.
func (am *AuthManager) SerializeSessionValueHandler() AuthManagerSerializeSessionValueHandler {
	return am.serializeSessionValueHandler
}

// WithParseSessionValueHandler sets the parse session value handler.
func (am *AuthManager) WithParseSessionValueHandler(handler AuthManagerParseSessionValueHandler) *AuthManager {
	am.parseSessionValueHandler = handler
	return am
}

// ParseSessionValueHandler returns the parse session value handler.
func (am *AuthManager) ParseSessionValueHandler() AuthManagerParseSessionValueHandler {
	return am.parseSessionValueHandler
}

// WithPersistHandler sets the persist handler.
func (am *AuthManager) WithPersistHandler(handler AuthManagerPersistHandler) *AuthManager {
	am.persistHandler = handler
	return am
}

// PersistHandler returns the persist handler.
func (am *AuthManager) PersistHandler() AuthManagerPersistHandler {
	return am.persistHandler
}

// WithFetchHandler sets the fetch handler.
func (am *AuthManager) WithFetchHandler(handler AuthManagerFetchHandler) *AuthManager {
	am.fetchHandler = handler
	return am
}

// FetchHandler returns the fetch handler.
// It is used in `VerifySession` to satisfy session cache misses.
func (am *AuthManager) FetchHandler() AuthManagerFetchHandler {
	return am.fetchHandler
}

// WithRemoveHandler sets the remove handler.
func (am *AuthManager) WithRemoveHandler(handler AuthManagerRemoveHandler) *AuthManager {
	am.removeHandler = handler
	return am
}

// RemoveHandler returns the remove handler.
// It is used in validate session if the session is found to be invalid.
func (am *AuthManager) RemoveHandler() AuthManagerRemoveHandler {
	return am.removeHandler
}

// WithValidateHandler sets the validate handler.
func (am *AuthManager) WithValidateHandler(handler AuthManagerValidateHandler) *AuthManager {
	am.validateHandler = handler
	return am
}

// ValidateHandler returns the validate handler.
func (am *AuthManager) ValidateHandler() AuthManagerValidateHandler {
	return am.validateHandler
}

// WithLoginRedirectHandler sets the login redirect handler.
func (am *AuthManager) WithLoginRedirectHandler(handler AuthManagerRedirectHandler) *AuthManager {
	am.loginRedirectHandler = handler
	return am
}

// LoginRedirectHandler returns the login redirect handler.
func (am *AuthManager) LoginRedirectHandler() AuthManagerRedirectHandler {
	return am.loginRedirectHandler
}

// WithPostLoginRedirectHandler sets the post login redirect handler.
func (am *AuthManager) WithPostLoginRedirectHandler(handler AuthManagerRedirectHandler) *AuthManager {
	am.postLoginRedirectHandler = handler
	return am
}

// PostLoginRedirectHandler returns the redirect handler for login complete.
func (am *AuthManager) PostLoginRedirectHandler() AuthManagerRedirectHandler {
	return am.postLoginRedirectHandler
}

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

func (am AuthManager) expire(ctx *Ctx, sessionValue string) error {
	ctx.ExpireCookie(am.CookieName(), am.CookiePath())
	// if we have a remove handler and the sessionID is set
	if am.removeHandler != nil {
		err := am.removeHandler(ctx.Context(), sessionValue, ctx.state)
		if err != nil {
			return err
		}
	}
	return nil
}

func (am AuthManager) shouldUpdateSessionExpiry() bool {
	return am.sessionTimeoutProvider != nil
}

// InjectCookie injects a session cookie into the context.
func (am *AuthManager) injectCookie(ctx *Ctx, name, value string, expire time.Time) {
	ctx.WriteNewCookie(name, value, expire, am.CookiePath(), am.CookiesHTTPSOnly())
}

// readParam reads a param from a given request context from either the cookies or headers.
func (am *AuthManager) readParam(name string, ctx *Ctx) (output string) {
	if cookie := ctx.GetCookie(name); cookie != nil {
		output = cookie.Value
	}
	return
}

// ReadSessionID reads a session id from a given request context.
func (am *AuthManager) readSessionValue(ctx *Ctx) string {
	return am.readParam(am.CookieName(), ctx)
}
