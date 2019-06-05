package web

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
)

// Option is an option for an app.
type Option func(*App)

// OptConfig sets the config.
func OptConfig(cfg Config) Option {
	return func(a *App) {
		a.Config = cfg
		a.Auth = NewAuthManager(cfg)
		a.Views = NewViewCache(OptViewCacheConfig(&cfg.Views))
	}
}

// OptConfigFromEnv sets the config from the environment.
func OptConfigFromEnv() Option {
	return func(a *App) {
		var cfg Config
		env.Env().ReadInto(&cfg)
		a.Config = cfg
		a.Auth = NewAuthManager(cfg)
		a.Views = NewViewCache(OptViewCacheConfig(&cfg.Views))
	}
}

// OptBindAddr sets the config bind address
func OptBindAddr(bindAddr string) Option {
	return func(a *App) {
		a.Config.BindAddr = bindAddr
	}
}

// OptPort sets the config bind address
func OptPort(port int32) Option {
	return func(a *App) {
		a.Config.Port = port
		a.Config.BindAddr = fmt.Sprintf(":%v", port)
	}
}

// OptLog sets the logger.
func OptLog(log logger.Log) Option {
	return func(a *App) { a.Log = log }
}

// OptServer sets the underlying server.
func OptServer(server *http.Server) Option {
	return func(a *App) { a.Server = server }
}

// OptAuth sets the auth manager.
func OptAuth(auth AuthManager) Option {
	return func(a *App) { a.Auth = auth }
}

// OptTracer sets the tracer.
func OptTracer(tracer Tracer) Option {
	return func(a *App) { a.Tracer = tracer }
}

// OptViews sets the view cache.
func OptViews(views *ViewCache) Option {
	return func(a *App) { a.Views = views }
}

// OptTLSConfig sets the tls config.
func OptTLSConfig(cfg *tls.Config) Option {
	return func(a *App) { a.TLSConfig = cfg }
}

// OptDefaultHeader sets a default header.
func OptDefaultHeader(key, value string) Option {
	return func(a *App) {
		if a.DefaultHeaders == nil {
			a.DefaultHeaders = make(http.Header)
		}
		a.DefaultHeaders.Set(key, value)
	}
}

// OptDefaultHeaders sets default headers.
func OptDefaultHeaders(headers http.Header) Option {
	return func(a *App) { a.DefaultHeaders = headers }
}

// OptDefaultMiddleware sets default middleware.
func OptDefaultMiddleware(middleware ...Middleware) Option {
	return func(a *App) { a.DefaultMiddleware = middleware }
}

// OptUse adds to the default middleware.
func OptUse(m Middleware) Option {
	return func(a *App) { a.DefaultMiddleware = append(a.DefaultMiddleware, m) }
}

// OptMethodNotAllowedHandler sets default headers.
func OptMethodNotAllowedHandler(action Action) Option {
	return func(a *App) { a.MethodNotAllowedHandler = a.RenderAction(action) }
}

// OptNotFoundHandler sets default headers.
func OptNotFoundHandler(action Action) Option {
	return func(a *App) { a.NotFoundHandler = a.RenderAction(action) }
}
