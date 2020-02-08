package web

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// Option is an option for an app.
type Option func(*App) error

// OptConfig sets the config.
func OptConfig(cfg Config) Option {
	return func(a *App) error {
		var err error
		a.Auth, err = NewAuthManager(OptAuthManagerFromConfig(cfg))
		if err != nil {
			return err
		}
		a.Config = cfg
		a.Views = NewViewCache(OptViewCacheConfig(&cfg.Views))
		return nil
	}
}

// OptConfigFromEnv sets the config from the environment.
func OptConfigFromEnv() Option {
	return func(a *App) error {
		var cfg Config
		if err := env.Env().ReadInto(&cfg); err != nil {
			return err
		}
		var err error
		a.Auth, err = NewAuthManager(OptAuthManagerFromConfig(cfg))
		if err != nil {
			return err
		}
		a.Config = cfg
		a.Views = NewViewCache(OptViewCacheConfig(&cfg.Views))
		return nil
	}
}

// OptBindAddr sets the config bind address
func OptBindAddr(bindAddr string) Option {
	return func(a *App) error {
		a.Config.BindAddr = bindAddr
		return nil
	}
}

// OptPort sets the config bind address
func OptPort(port int32) Option {
	return func(a *App) error {
		a.Config.Port = port
		a.Config.BindAddr = fmt.Sprintf(":%v", port)
		return nil
	}
}

// OptLog sets the logger.
func OptLog(log logger.Log) Option {
	return func(a *App) error {
		a.Log = log
		return nil
	}
}

// OptServer sets the underlying server.
func OptServer(server *http.Server) Option {
	return func(a *App) error {
		a.Server = server
		return nil
	}
}

// OptAuth sets the auth manager.
func OptAuth(auth AuthManager, err error) Option {
	return func(a *App) error {
		if err != nil {
			return err
		}
		a.Auth = auth
		return nil
	}
}

// OptTracer sets the tracer.
func OptTracer(tracer Tracer) Option {
	return func(a *App) error {
		a.Tracer = tracer
		return nil
	}
}

// OptViews sets the view cache.
func OptViews(views *ViewCache) Option {
	return func(a *App) error {
		a.Views = views
		return nil
	}
}

// OptTLSConfig sets the tls config.
func OptTLSConfig(cfg *tls.Config) Option {
	return func(a *App) error {
		a.TLSConfig = cfg
		return nil
	}
}

// OptDefaultHeader sets a default header.
func OptDefaultHeader(key, value string) Option {
	return func(a *App) error {
		if a.DefaultHeaders == nil {
			a.DefaultHeaders = make(http.Header)
		}
		a.DefaultHeaders.Set(key, value)
		return nil
	}
}

// OptDefaultHeaders sets default headers.
func OptDefaultHeaders(headers http.Header) Option {
	return func(a *App) error {
		a.DefaultHeaders = headers
		return nil
	}
}

// OptDefaultMiddleware sets default middleware.
func OptDefaultMiddleware(middleware ...Middleware) Option {
	return func(a *App) error {
		a.DefaultMiddleware = middleware
		return nil
	}
}

// OptUse adds to the default middleware.
func OptUse(m Middleware) Option {
	return func(a *App) error {
		a.DefaultMiddleware = append(a.DefaultMiddleware, m)
		return nil
	}
}

// OptMethodNotAllowedHandler sets default headers.
func OptMethodNotAllowedHandler(action Action) Option {
	return func(a *App) error {
		a.MethodNotAllowedHandler = a.RenderAction(action)
		return nil
	}
}

// OptNotFoundHandler sets default headers.
func OptNotFoundHandler(action Action) Option {
	return func(a *App) error {
		a.NotFoundHandler = a.RenderAction(action)
		return nil
	}
}

// OptShutdownGracePeriod sets the shutdown grace period.
func OptShutdownGracePeriod(d time.Duration) Option {
	return func(a *App) error {
		a.Config.ShutdownGracePeriod = d
		return nil
	}
}

// OptServerOptions applies options to the underlying http server.
func OptServerOptions(opts ...webutil.HTTPServerOption) Option {
	return func(a *App) error {
		a.ServerOptions = opts
		return nil
	}
}
