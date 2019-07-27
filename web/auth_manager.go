package web

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/ex"
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
	for _, opt := range options {
		if err = opt(&manager); err != nil {
			return
		}
	}
	return
}

// AuthManagerOption is a variadic option for auth managers.
type AuthManagerOption func(*AuthManager) error

// OptAuthManagerFromConfig returns an auth manager from a config.
func OptAuthManagerFromConfig(cfg Config) AuthManagerOption {
	return func(am *AuthManager) error {
		switch cfg.AuthManagerModeOrDefault() {
		case AuthManagerModeJWT:
			*am = NewJWTAuthManager(cfg.MustAuthSecret())
		case AuthManagerModeLocal: // local should only be used for debugging.
			*am = NewLocalAuthManager()
		case AuthManagerModeRemote:
			*am = NewRemoteAuthManager()
		default:
			return ex.New("invalid auth manager mode")
		}

		am.CookieSecure = cfg.CookieSecureOrDefault()
		am.CookieHTTPOnly = cfg.CookieHTTPOnlyOrDefault()
		am.CookieName = cfg.CookieNameOrDefault()
		am.CookiePath = cfg.CookiePathOrDefault()
		am.CookieSameSite = cfg.CookieSameSiteOrDefault()
		am.SessionTimeoutProvider = SessionTimeoutProvider(!cfg.SessionTimeoutIsRelative, cfg.SessionTimeoutOrDefault())
		return nil
	}
}

// NewRemoteAuthManager returns an empty auth manager.
// You must provide a fetch, persist, and remove handler, and optionally a login redirect handler.
func NewRemoteAuthManager() AuthManager {
	return AuthManager{
		Mode:           AuthManagerModeRemote,
		CookieSecure:   DefaultCookieSecure,
		CookieHTTPOnly: DefaultCookieHTTPOnly,
		CookieSameSite: DefaultCookieSameSite,
	}
}

// NewLocalAuthManagerFromCache returns a new locally cached session manager that saves sessions to the cache provided
func NewLocalAuthManagerFromCache(cache *LocalSessionCache) AuthManager {
	return AuthManager{
		Mode:           AuthManagerModeLocal,
		CookieSecure:   DefaultCookieSecure,
		CookieHTTPOnly: DefaultCookieHTTPOnly,
		CookieSameSite: DefaultCookieSameSite,
		PersistHandler: cache.PersistHandler,
		FetchHandler:   cache.FetchHandler,
		RemoveHandler:  cache.RemoveHandler,
	}
}

// NewLocalAuthManager returns a new locally cached session manager.
// It saves sessions to a local store.
func NewLocalAuthManager() AuthManager {
	cache := NewLocalSessionCache()
	return AuthManager{
		Mode:           AuthManagerModeJWT,
		CookieSecure:   DefaultCookieSecure,
		CookieHTTPOnly: DefaultCookieHTTPOnly,
		CookieSameSite: DefaultCookieSameSite,
		PersistHandler: cache.PersistHandler,
		FetchHandler:   cache.FetchHandler,
		RemoveHandler:  cache.RemoveHandler,
	}
}

// NewJWTAuthManager returns a new jwt session manager.
// It issues JWT tokens to identify users.
func NewJWTAuthManager(key []byte) AuthManager {
	jwtm := NewJWTManager(key)
	return AuthManager{
		SerializeSessionValueHandler: jwtm.SerializeSessionValueHandler,
		ParseSessionValueHandler:     jwtm.ParseSessionValueHandler,
		SessionTimeoutProvider:       SessionTimeoutProviderAbsolute(DefaultSessionTimeout),
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
	// Mode is the mechanism the auth manager tracks sessions.
	// Possible values include local, remote, jwt.
	Mode AuthManagerMode

	CookieName     string
	CookiePath     string
	CookieSecure   bool
	CookieHTTPOnly bool
	CookieSameSite string

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
	am.injectCookie(ctx, am.CookieNameOrDefault(), sessionValue, session.ExpiresUTC)
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
	ctx.ExpireCookie(am.CookieNameOrDefault(), am.CookiePathOrDefault())
	ctx.Session = nil

	// call the remove handler if one has been provided
	if am.RemoveHandler != nil {
		return am.RemoveHandler(ctx.Context(), sessionValue)
	}
	return nil
}

// VerifySession checks a sessionID to see if it's valid.
// It also handles updating a rolling expiry.
func (am AuthManager) VerifySession(ctx *Ctx) (session *Session, err error) {
	// pull the sessionID off the request
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
				am.expire(ctx, sessionValue)
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
		am.injectCookie(ctx, am.CookieNameOrDefault(), sessionValue, session.ExpiresUTC)
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
	return ctx.App.DefaultProvider.NotAuthorized()
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

// CookieNameOrDefault returns the cookie name or a default.
func (am AuthManager) CookieNameOrDefault() string {
	if am.CookieName == "" {
		return DefaultCookieName
	}
	return am.CookieName
}

// CookiePathOrDefault returns the session param path.
func (am AuthManager) CookiePathOrDefault() string {
	if am.CookiePath == "" {
		return DefaultCookiePath
	}
	return am.CookiePath
}

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

func (am AuthManager) expire(ctx *Ctx, sessionValue string) error {
	ctx.ExpireCookie(am.CookieNameOrDefault(), am.CookiePathOrDefault())
	// if we have a remove handler and the sessionID is set
	if am.RemoveHandler != nil {
		err := am.RemoveHandler(ctx.Context(), sessionValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (am AuthManager) shouldUpdateSessionExpiry() bool {
	return am.SessionTimeoutProvider != nil
}

// InjectCookie injects a session cookie into the context.
func (am AuthManager) injectCookie(ctx *Ctx, name, value string, expire time.Time) {
	ctx.WriteNewCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		Path:     am.CookiePathOrDefault(),
		HttpOnly: am.CookieHTTPOnly,
		Secure:   am.CookieSecure,
		SameSite: webutil.MustParseSameSite(am.CookieSameSite),
	})
}

// readParam reads a param from a given request context from either the cookies or headers.
func (am AuthManager) readParam(name string, ctx *Ctx) (output string) {
	if cookie := ctx.Cookie(name); cookie != nil {
		output = cookie.Value
	}
	return
}

// ReadSessionID reads a session id from a given request context.
func (am AuthManager) readSessionValue(ctx *Ctx) string {
	return am.readParam(am.CookieNameOrDefault(), ctx)
}
