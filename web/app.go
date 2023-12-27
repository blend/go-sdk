/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/proxyprotocol"
	"github.com/blend/go-sdk/webutil"
)

// MustNew creates a new app and panics if there is an error.
func MustNew(options ...Option) *App {
	app, err := New(options...)
	if err != nil {
		panic(err)
	}
	return app
}

// New returns a new web app.
func New(options ...Option) (*App, error) {
	views, err := NewViewCache()
	if err != nil {
		return nil, err
	}
	auth, err := NewAuthManager()
	if err != nil {
		return nil, err
	}
	a := App{
		RouteTree:		new(RouteTree),
		Auth:			auth,
		BaseContext:		func(_ net.Listener) context.Context { return context.Background() },
		BaseHeaders:		BaseHeaders(),
		BaseState:		new(SyncState),
		DefaultProvider:	views,
		Latch:			async.NewLatch(),
		Server:			new(http.Server),
		Statics:		map[string]*StaticFileServer{},
		Views:			views,
	}

	for _, option := range options {
		if err = option(&a); err != nil {
			return nil, err
		}
	}
	return &a, nil
}

// App is the server for the app.
type App struct {
	*async.Latch
	*RouteTree

	Config	Config

	Auth		AuthManager
	BaseContext	func(net.Listener) context.Context

	BaseHeaders	http.Header
	BaseMiddleware	[]Middleware
	BaseState	State

	Log	logger.Log
	Tracer	Tracer

	TLSConfig	*tls.Config
	Server		*http.Server
	Listener	net.Listener

	Statics	map[string]*StaticFileServer

	DefaultProvider	ResultProvider
	Views		*ViewCache

	PanicAction	PanicAction
}

// Background returns a base context.
func (a *App) Background() context.Context {
	if a.BaseContext != nil {
		return a.BaseContext(a.Listener)
	}
	return context.Background()
}

// --------------------------------------------------------------------------------
// Lifecycle
// --------------------------------------------------------------------------------

// Start starts the server and binds to the given address.
func (a *App) Start() (err error) {
	if !a.Latch.CanStart() {
		return ex.New(async.ErrCannotStart)
	}
	for _, opt := range a.httpServerOptions() {
		if err = opt(a.Server); err != nil {
			return err
		}
	}

	err = a.StartupTasks()
	if err != nil {
		return
	}

	var shutdownErr error
	if a.Listener == nil {
		serverProtocol := "http"
		if a.Server.TLSConfig != nil {
			serverProtocol = "https (tls)"
		}
		if a.Server.Addr == "" {
			a.Server.Addr = a.Config.BindAddrOrDefault()
		}

		var rawListener net.Listener
		rawListener, err = net.Listen("tcp", a.Server.Addr)
		if err != nil {
			err = ex.New(err)
			return
		}
		typedListener, ok := rawListener.(*net.TCPListener)
		if !ok {
			err = ex.New("listener returned was not a net.TCPListener")
			return
		}
		a.Listener = webutil.TCPKeepAliveListener{
			TCPListener:		typedListener,
			KeepAlive:		a.Config.KeepAliveOrDefault(),
			KeepAlivePeriod:	a.Config.KeepAlivePeriodOrDefault(),
		}

		if a.Config.UseProxyProtocol {
			logger.MaybeInfofContext(a.Background(), a.Log, "%s using proxy protocol", serverProtocol)
			a.Listener = &proxyprotocol.Listener{Listener: a.Listener}
		}

		if a.Server.TLSConfig != nil {
			logger.MaybeInfofContext(a.Background(), a.Log, "%s using tls", serverProtocol)
			a.Listener = tls.NewListener(a.Listener, a.Server.TLSConfig)
		}
		logger.MaybeInfofContext(a.Background(), a.Log, "%s server started, listening on %s", serverProtocol, a.Server.Addr)
	} else {
		logger.MaybeInfofContext(a.Background(), a.Log, "http server started, using custom listener")
	}

	a.Started()
	shutdownErr = a.Server.Serve(a.Listener)
	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		err = ex.New(shutdownErr)
	}
	logger.MaybeInfofContext(a.Background(), a.Log, "server stopped serving")
	a.Stopped()
	return
}

// Stop stops the server.
func (a *App) Stop() error {
	if !a.CanStop() {
		return ex.New(async.ErrCannotStop)
	}
	a.Stopping()

	ctx := a.Background()
	var cancel context.CancelFunc
	if gracePeriod := a.Config.ShutdownGracePeriodOrDefault(); gracePeriod > 0 {
		logger.MaybeInfofContext(ctx, a.Log, "server shutdown grace period: %v", gracePeriod)
		ctx, cancel = context.WithTimeout(ctx, gracePeriod)
		defer cancel()
	}
	logger.MaybeInfofContext(ctx, a.Log, "server keep alives disabled")
	a.Server.SetKeepAlivesEnabled(false)
	logger.MaybeInfofContext(ctx, a.Log, "server shutting down")
	if err := a.Server.Shutdown(ctx); err != nil {
		if err == context.DeadlineExceeded {
			logger.MaybeWarningfContext(ctx, a.Log, "server shutdown grace period exceeded, connections forcibly closed")
		} else {
			return ex.New(err)
		}
	}
	logger.MaybeInfofContext(a.Background(), a.Log, "server shutdown complete")
	return nil
}

// --------------------------------------------------------------------------------
// Register Controllers
// --------------------------------------------------------------------------------

// Register registers controllers with the app's router.
func (a *App) Register(controllers ...Controller) {
	for _, c := range controllers {
		c.Register(a)
	}
}

// --------------------------------------------------------------------------------
// Static Result Methods
// --------------------------------------------------------------------------------

// ServeStatic serves files from the given file system root(s)..
// If the path does not end with "/*filepath" that suffix will be added for you internally.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
func (a *App) ServeStatic(route string, searchPaths []string, middleware ...Middleware) {
	var searchPathFS []http.FileSystem
	for _, searchPath := range searchPaths {
		searchPathFS = append(searchPathFS, http.Dir(searchPath))
	}
	sfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(searchPathFS...),
		OptStaticFileServerCacheDisabled(true),
	)
	mountedRoute := a.formatStaticMountRoute(route)
	a.Statics[mountedRoute] = sfs
	a.Method(webutil.MethodGet, mountedRoute, sfs.Action, middleware...)
}

// ServeStaticCached serves files from the given file system root(s).
// If the path does not end with "/*filepath" that suffix will be added for you internally.
func (a *App) ServeStaticCached(route string, searchPaths []string, middleware ...Middleware) {
	var searchPathFS []http.FileSystem
	for _, searchPath := range searchPaths {
		searchPathFS = append(searchPathFS, http.Dir(searchPath))
	}
	sfs := NewStaticFileServer(
		OptStaticFileServerSearchPaths(searchPathFS...),
	)
	mountedRoute := a.formatStaticMountRoute(route)
	a.Statics[mountedRoute] = sfs
	a.Method(webutil.MethodGet, mountedRoute, sfs.Action, middleware...)
}

// SetStaticRewriteRule adds a rewrite rule for a specific statically served path.
// It mutates the path for the incoming static file request to the fileserver according to the action.
func (a *App) SetStaticRewriteRule(route, match string, action RewriteAction) error {
	mountedRoute := a.formatStaticMountRoute(route)
	if static, hasRoute := a.Statics[mountedRoute]; hasRoute {
		return static.AddRewriteRule(match, action)
	}
	return ex.New("no static fileserver mounted at route", ex.OptMessagef("route: %s", route))
}

// SetStaticHeader adds a header for the given static path.
// These headers are automatically added to any result that the static path fileserver sends.
func (a *App) SetStaticHeader(route, key, value string) error {
	mountedRoute := a.formatStaticMountRoute(route)
	if static, hasRoute := a.Statics[mountedRoute]; hasRoute {
		static.AddHeader(key, value)
		return nil
	}
	return ex.New("no static fileserver mounted at route", ex.OptMessagef("route: %s", mountedRoute))
}

// --------------------------------------------------------------------------------
// Route Registration / HTTP Methods
// --------------------------------------------------------------------------------

// GET registers a GET request route handler with the given middleware.
func (a *App) GET(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodGet, path, action, middleware...)
}

// OPTIONS registers a OPTIONS request route handler the given middleware.
func (a *App) OPTIONS(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodOptions, path, action, middleware...)
}

// HEAD registers a HEAD request route handler with the given middleware.
func (a *App) HEAD(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodHead, path, action, middleware...)
}

// PUT registers a PUT request route handler with the given middleware.
func (a *App) PUT(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodPut, path, action, middleware...)
}

// PATCH registers a PATCH request route handler with the given middleware.
func (a *App) PATCH(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodPatch, path, action, middleware...)
}

// POST registers a POST request route handler with the given middleware.
func (a *App) POST(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodPost, path, action, middleware...)
}

// DELETE registers a DELETE request route handler with the given middleware.
func (a *App) DELETE(path string, action Action, middleware ...Middleware) {
	a.Method(http.MethodDelete, path, action, middleware...)
}

// Method registers an action for a given method and path with the given middleware.
func (a *App) Method(method string, path string, action Action, middleware ...Middleware) {
	a.RouteTree.Handle(method, path, a.RenderAction(NestMiddleware(action, append(middleware, a.BaseMiddleware...)...)))
}

// MethodBare registers an action for a given method and path with the given middleware that omits logging and tracing.
func (a *App) MethodBare(method string, path string, action Action, middleware ...Middleware) {
	a.RouteTree.Handle(method, path, a.RenderActionBare(NestMiddleware(action, append(middleware, a.BaseMiddleware...)...)))
}

// Lookup finds the route data for a given method and path.
func (a *App) Lookup(method, path string) (route *Route, params RouteParameters, skipSlashRedirect bool) {
	if root := a.RouteTree.Routes[method]; root != nil {
		route, params, skipSlashRedirect = root.getValue(path)
		return
	}
	return
}

// --------------------------------------------------------------------------------
// Request Pipeline
// --------------------------------------------------------------------------------

// ServeHTTP makes the router implement the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !a.Config.DisablePanicRecovery {
		defer a.recover(w, req)
	}
	// load the request start time onto the request.
	req = req.WithContext(WithRequestStarted(req.Context(), time.Now().UTC()))
	a.RouteTree.ServeHTTP(w, req)
}

// RenderAction is the translation step from Action to Handler.
func (a *App) RenderAction(action Action) Handler {
	return func(w http.ResponseWriter, r *http.Request, route *Route, p RouteParameters) {
		ctx := NewCtx(webutil.NewStatusResponseWriter(w), r, a.ctxOptions(r.Context(), route, p)...)
		defer ctx.Close()
		defer a.logRequest(ctx)

		var err error
		if a.Tracer != nil {
			tf := ctx.Tracer.Start(ctx)
			defer func() {
				tf.Finish(ctx, err)
			}()
		}

		for key, value := range a.BaseHeaders {
			ctx.Response.Header()[key] = value
		}

		if result := action(ctx); result != nil {
			if typed, ok := result.(ResultPreRender); ok {
				if errPreRender := typed.PreRender(ctx); errPreRender != nil {
					a.maybeLogFatal(ctx.Context(), errPreRender, ctx.Request)
					err = ex.New(errPreRender, ex.OptInner(err))
				}
			}
			if errRender := result.Render(ctx); errRender != nil {
				a.maybeLogFatal(ctx.Context(), errRender, ctx.Request)
				err = ex.New(errRender, ex.OptInner(err))
			}
			if typed, ok := result.(ResultPostRender); ok {
				if errPostRender := typed.PostRender(ctx); errPostRender != nil {
					a.maybeLogFatal(ctx.Context(), errPostRender, ctx.Request)
					err = ex.New(errPostRender, ex.OptInner(err))
				}
			}
		}
	}
}

// RenderActionBare is the translation step from Action to Handler that omits logging.
func (a *App) RenderActionBare(action Action) Handler {
	return func(w http.ResponseWriter, r *http.Request, route *Route, p RouteParameters) {
		ctx := NewCtx(webutil.NewStatusResponseWriter(w), r, a.ctxOptions(r.Context(), route, p)...)
		defer ctx.Close()

		for key, value := range a.BaseHeaders {
			ctx.Response.Header()[key] = value
		}

		if result := action(ctx); result != nil {
			if typed, ok := result.(ResultPreRender); ok {
				if errPreRender := typed.PreRender(ctx); errPreRender != nil {
					a.maybeLogFatal(ctx.Context(), errPreRender, ctx.Request)
				}
			}
			if err := result.Render(ctx); err != nil {
				a.maybeLogFatal(ctx.Context(), err, ctx.Request)
			}
			if typed, ok := result.(ResultPostRender); ok {
				if errPostRender := typed.PostRender(ctx); errPostRender != nil {
					a.maybeLogFatal(ctx.Context(), errPostRender, ctx.Request)
				}
			}
		}
	}
}

//
// startup helpers
//

// StartupTasks runs common startup tasks.
// These tasks include anything outside setting up the underlying server itself.
// Right now, this is limited to initializing the view cache if relevant.
func (a *App) StartupTasks() (err error) {
	if err = a.Views.Initialize(); err != nil {
		return
	}
	return nil
}

//
// internal helpers
//

func (a *App) formatStaticMountRoute(route string) string {
	mountedRoute := route
	if !strings.HasSuffix(mountedRoute, "*"+RouteTokenFilepath) {
		if strings.HasSuffix(mountedRoute, "/") {
			mountedRoute = mountedRoute + "*" + RouteTokenFilepath
		} else {
			mountedRoute = mountedRoute + "/*" + RouteTokenFilepath
		}
	}
	return mountedRoute
}

func (a *App) httpServerOptions() []webutil.HTTPServerOption {
	return []webutil.HTTPServerOption{
		webutil.OptHTTPServerHandler(a),
		webutil.OptHTTPServerTLSConfig(a.TLSConfig),
		webutil.OptHTTPServerAddr(a.Config.BindAddrOrDefault()),
		webutil.OptHTTPServerMaxHeaderBytes(a.Config.MaxHeaderBytesOrDefault()),
		webutil.OptHTTPServerReadTimeout(a.Config.ReadTimeoutOrDefault()),
		webutil.OptHTTPServerReadHeaderTimeout(a.Config.ReadHeaderTimeoutOrDefault()),
		webutil.OptHTTPServerWriteTimeout(a.Config.WriteTimeoutOrDefault()),
		webutil.OptHTTPServerIdleTimeout(a.Config.IdleTimeoutOrDefault()),
		webutil.OptHTTPServerBaseContext(a.BaseContext),
	}
}

func (a *App) ctxOptions(ctx context.Context, route *Route, p RouteParameters) []CtxOption {
	return []CtxOption{
		OptCtxApp(a),
		OptCtxAuth(a.Auth),
		OptCtxDefaultProvider(a.DefaultProvider),
		OptCtxViews(a.Views),
		OptCtxRoute(route),
		OptCtxRouteParams(p),
		OptCtxState(a.BaseState.Copy()),
		OptCtxLog(a.Log),
		OptCtxTracer(a.Tracer),
		OptCtxRequestStarted(GetRequestStarted(ctx)),
	}
}

func (a *App) recover(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		err := ex.New(rcv)
		a.maybeLogFatal(req.Context(), err, req)
		if a.PanicAction != nil {
			a.RenderAction(func(ctx *Ctx) Result {
				return a.PanicAction(ctx, err)
			})(w, req, nil, nil)
			return
		}
		http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
		return
	}
}

func (a *App) maybeLogFatal(ctx context.Context, err error, req *http.Request) {
	if !logger.IsLoggerSet(a.Log) || err == nil {
		return
	}
	a.Log.TriggerContext(
		ctx,
		logger.NewErrorEvent(
			logger.Fatal,
			err,
			logger.OptErrorEventState(req),
		),
	)
}

func (a *App) logRequest(r *Ctx) {
	requestEvent := webutil.NewHTTPRequestEvent(r.Request.Clone(r.Context()),
		webutil.OptHTTPRequestStatusCode(r.Response.StatusCode()),
		webutil.OptHTTPRequestContentLength(r.Response.ContentLength()),
		webutil.OptHTTPRequestHeader(r.Response.Header().Clone()),
		webutil.OptHTTPRequestElapsed(r.Elapsed()),
	)
	if r.Route != nil {
		requestEvent.Route = r.Route.String()
	}
	if requestEvent.Header != nil {
		requestEvent.ContentType = requestEvent.Header.Get(webutil.HeaderContentType)
		requestEvent.ContentEncoding = requestEvent.Header.Get(webutil.HeaderContentEncoding)
	}
	a.maybeLogTrigger(r.Context(), r.Log, requestEvent)
}

func (a *App) maybeLogTrigger(ctx context.Context, log logger.Log, e logger.Event) {
	if !logger.IsLoggerSet(log) || e == nil {
		return
	}
	log.TriggerContext(ctx, e)
}
