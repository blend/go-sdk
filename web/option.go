/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
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
		a.BaseHeaders = MergeHeaders(BaseHeaders(), CopySingleHeaders(cfg.DefaultHeaders))
		a.Views, err = NewViewCache(OptViewCacheConfig(&cfg.Views))
		return err
	}
}

// OptConfigFromEnv sets the config from the environment.
func OptConfigFromEnv() Option {
	return func(a *App) error {
		var cfg Config
		if err := env.Env().ReadInto(&cfg); err != nil {
			return err
		}
		return OptConfig(cfg)(a)
	}
}

// OptBindAddr sets the config bind address.
func OptBindAddr(bindAddr string) Option {
	return func(a *App) error {
		a.Config.BindAddr = bindAddr
		return nil
	}
}

// OptPort sets the config bind address.
func OptPort(port int32) Option {
	return func(a *App) error {
		a.Config.Port = port
		a.Config.BindAddr = fmt.Sprintf(":%v", port)
		return nil
	}
}

// OptBaseContext sets a base context on an `App`; this is a context that
// is always used via the `BaseContext` function, independent of the listener.
func OptBaseContext(ctx context.Context) Option {
	return OptBaseContextFunc(func(_ net.Listener) context.Context {
		return ctx
	})
}

// OptBaseContextFunc sets a base context function on an `App`; this
// context function will propagated from the `App` to `http.Server.BaseContext`.
func OptBaseContextFunc(bc func(net.Listener) context.Context) Option {
	return func(a *App) error {
		a.BaseContext = bc
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
		if a.BaseHeaders == nil {
			a.BaseHeaders = make(http.Header)
		}
		a.BaseHeaders.Set(key, value)
		return nil
	}
}

// OptDefaultHeaders sets base headers.
//
// DEPRECATION(1.2021*): this method will be removed.
func OptDefaultHeaders(headers http.Header) Option {
	return func(a *App) error {
		a.BaseHeaders = headers
		return nil
	}
}

// OptBaseHeaders sets base headers.
func OptBaseHeaders(headers http.Header) Option {
	return func(a *App) error {
		a.BaseHeaders = headers
		return nil
	}
}

// OptDefaultMiddleware sets base middleware.
//
// DEPRECATION(1.2021*): this method will be removed.
func OptDefaultMiddleware(middleware ...Middleware) Option {
	return func(a *App) error {
		a.BaseMiddleware = middleware
		return nil
	}
}

// OptBaseMiddleware sets default middleware.
func OptBaseMiddleware(middleware ...Middleware) Option {
	return func(a *App) error {
		a.BaseMiddleware = middleware
		return nil
	}
}

// OptUse adds to the default middleware.
func OptUse(m Middleware) Option {
	return func(a *App) error {
		a.BaseMiddleware = append(a.BaseMiddleware, m)
		return nil
	}
}

// OptBaseStateValue sets a base state value.
func OptBaseStateValue(key string, value interface{}) Option {
	return func(a *App) error {
		a.BaseState.Set(key, value)
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
//
// Many of the fields on the server are overwritten by the config on `app.Start`.
// You should only use `OptServerOptions` for fields that are not governed by the config
// such as the stdlib logger.
func OptServerOptions(opts ...webutil.HTTPServerOption) Option {
	return func(a *App) error {
		for _, opt := range opts {
			if err := opt(a.Server); err != nil {
				return err
			}
		}
		return nil
	}
}

// OptReadTimeout sets the read timeout.
//
// Note that this will override the config setting if OptConfig comes before it
// and will be overwritten by the config if OptConfig comes after it.
func OptReadTimeout(d time.Duration) Option {
	return func(a *App) error {
		a.Config.ReadTimeout = d
		return nil
	}
}

// OptReadHeaderTimeout sets the read header timeout.
//
// Note that this will override the config setting if OptConfig comes before it
// and will be overwritten by the config if OptConfig comes after it.
func OptReadHeaderTimeout(d time.Duration) Option {
	return func(a *App) error {
		a.Config.ReadHeaderTimeout = d
		return nil
	}
}

// OptWriteTimeout sets the write timeout.
//
// Note that this will override the config setting if OptConfig comes before it
// and will be overwritten by the config if OptConfig comes after it.
func OptWriteTimeout(d time.Duration) Option {
	return func(a *App) error {
		a.Config.WriteTimeout = d
		return nil
	}
}

// OptIdleTimeout sets the idle timeout.
//
// Note that this will override the config setting if OptConfig comes before it
// and will be overwritten by the config if OptConfig comes after it.
func OptIdleTimeout(d time.Duration) Option {
	return func(a *App) error {
		a.Config.IdleTimeout = d
		return nil
	}
}

// OptMaxHeaderBytes sets the max header bytes.
//
// Note that this will override the config setting if OptConfig comes before it
// and will be overwritten by the config if OptConfig comes after it.
func OptMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(a *App) error {
		a.Config.MaxHeaderBytes = maxHeaderBytes
		return nil
	}
}

// OptBaseURL sets the config base url.
func OptBaseURL(baseURL string) Option {
	return func(a *App) error {
		a.Config.BaseURL = baseURL
		return nil
	}
}
