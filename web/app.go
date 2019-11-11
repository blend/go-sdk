package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// New returns a new app.
func New() *App {
	views := NewViewCache()
	return &App{
		latch:                 async.NewLatch(),
		hsts:                  &HSTSConfig{},
		auth:                  &AuthManager{},
		bindAddr:              DefaultBindAddr,
		state:                 &SyncState{},
		statics:               map[string]Fileserver{},
		readTimeout:           DefaultReadTimeout,
		writeTimeout:          DefaultWriteTimeout,
		redirectTrailingSlash: true,
		recoverPanics:         true,
		defaultHeaders:        DefaultHeaders,
		shutdownGracePeriod:   DefaultShutdownGracePeriod,
		views:                 views,
		defaultResultProvider: views,
	}
}

// NewFromEnv returns a new app from the environment.
func NewFromEnv() (*App, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg), nil
}

// MustNewFromEnv returns a new app with a config set from environment
// variabales, and it will panic if there is an error.
func MustNewFromEnv() *App {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return NewFromConfig(cfg)
}

// NewFromConfig returns a new app from a given config.
func NewFromConfig(cfg *Config) *App {
	return New().WithConfig(cfg)
}

// App is the server for the app.
type App struct {
	latch *async.Latch
	cfg   *Config
	hsts  *HSTSConfig

	log   *logger.Logger
	auth  *AuthManager
	views *ViewCache

	baseURL  *url.URL
	bindAddr string

	tls      *tls.Config
	server   *http.Server
	handler  http.Handler
	listener *net.TCPListener

	// defaultHeaders are the default headers we apply to any request responses.
	defaultHeaders map[string]string
	// statics serve files at various routes
	statics map[string]Fileserver

	routes                  map[string]*node
	notFoundHandler         Handler
	methodNotAllowedHandler Handler
	panicAction             PanicAction
	redirectTrailingSlash   bool
	handleOptions           bool
	handleMethodNotAllowed  bool

	defaultMiddleware []Middleware
	tracer            Tracer

	defaultResultProvider ResultProvider

	maxHeaderBytes      int
	readTimeout         time.Duration
	readHeaderTimeout   time.Duration
	writeTimeout        time.Duration
	idleTimeout         time.Duration
	shutdownGracePeriod time.Duration

	state         *SyncState
	recoverPanics bool
}

// Latch returns the app lifecycle latch.
func (a *App) Latch() *async.Latch {
	return a.latch
}

// NotifyStarted returns the notify started chan.
func (a *App) NotifyStarted() <-chan struct{} {
	return a.latch.NotifyStarted()
}

// NotifyStopped returns the notify stopped chan.
func (a *App) NotifyStopped() <-chan struct{} {
	return a.latch.NotifyStopped()
}

// IsRunning returns if the app is running.
func (a *App) IsRunning() bool {
	return a.latch.IsRunning()
}

// WithConfig sets the config and applies the config's setting.
func (a *App) WithConfig(cfg *Config) *App {
	a.cfg = cfg

	a.WithBindAddr(cfg.GetBindAddr())
	a.WithRedirectTrailingSlash(cfg.GetRedirectTrailingSlash())
	a.WithHandleMethodNotAllowed(cfg.GetHandleMethodNotAllowed())
	a.WithHandleOptions(cfg.GetHandleOptions())
	a.WithRecoverPanics(cfg.GetRecoverPanics())
	a.WithDefaultHeaders(cfg.GetDefaultHeaders(DefaultHeaders))

	a.WithMaxHeaderBytes(cfg.GetMaxHeaderBytes())
	a.WithReadHeaderTimeout(cfg.GetReadHeaderTimeout())
	a.WithReadTimeout(cfg.GetReadTimeout())
	a.WithWriteTimeout(cfg.GetWriteTimeout())
	a.WithIdleTimeout(cfg.GetIdleTimeout())

	a.WithAuth(NewAuthManagerFromConfig(cfg))
	a.WithViews(NewViewCacheFromConfig(&cfg.Views))
	a.WithDefaultResultProvider(a.Views())
	a.WithBaseURL(webutil.MustParseURL(cfg.GetBaseURL()))
	a.WithShutdownGracePeriod(cfg.GetShutdownGracePeriod())

	a.WithHSTS(&cfg.HSTS)
	return a
}

// Config returns the app config.
func (a *App) Config() *Config {
	return a.cfg
}

// WithShutdownGracePeriod sets the shutdown grace period.
func (a *App) WithShutdownGracePeriod(gracePeriod time.Duration) *App {
	a.shutdownGracePeriod = gracePeriod
	return a
}

// ShutdownGracePeriod is the grace period on shutdown.
func (a *App) ShutdownGracePeriod() time.Duration {
	return a.shutdownGracePeriod
}

// WithDefaultHeaders sets the default headers
func (a *App) WithDefaultHeaders(headers map[string]string) *App {
	a.defaultHeaders = headers
	return a
}

// WithDefaultHeader adds a default header.
func (a *App) WithDefaultHeader(key string, value string) *App {
	a.defaultHeaders[key] = value
	return a
}

// DefaultHeaders returns the default headers.
func (a *App) DefaultHeaders() map[string]string {
	return a.defaultHeaders
}

// WithStateValue sets app state and returns a reference to the app for building apps with a fluent api.
func (a *App) WithStateValue(key string, value interface{}) *App {
	a.state.Set(key, value)
	return a
}

// StateValue gets app state element by key.
func (a *App) StateValue(key string) interface{} {
	return a.state.Get(key)
}

// State is a bag for common app state.
func (a *App) State() State {
	return a.state
}

// WithRedirectTrailingSlash sets if we should redirect missing trailing slashes.
func (a *App) WithRedirectTrailingSlash(value bool) *App {
	a.redirectTrailingSlash = value
	return a
}

// RedirectTrailingSlash returns if we should redirect missing trailing slashes to the correct route.
func (a *App) RedirectTrailingSlash() bool {
	return a.redirectTrailingSlash
}

// WithHandleMethodNotAllowed sets if we should handlem ethod not allowed.
func (a *App) WithHandleMethodNotAllowed(handle bool) *App {
	a.handleMethodNotAllowed = handle
	return a
}

// HandleMethodNotAllowed returns if we should handle unhandled verbs.
func (a *App) HandleMethodNotAllowed() bool {
	return a.handleMethodNotAllowed
}

// WithHandleOptions returns if we should handle OPTIONS requests.
func (a *App) WithHandleOptions(handle bool) *App {
	a.handleOptions = handle
	return a
}

// HandleOptions returns if we should handle OPTIONS requests.
func (a *App) HandleOptions() bool {
	return a.handleOptions
}

// WithRecoverPanics sets if the app should recover panics.
func (a *App) WithRecoverPanics(value bool) *App {
	a.recoverPanics = value
	return a
}

// RecoverPanics returns if the app recovers panics.
func (a *App) RecoverPanics() bool {
	return a.recoverPanics
}

// WithBaseURL sets the `BaseURL` field and returns a reference to the app for building apps with a fluent api.
func (a *App) WithBaseURL(baseURL *url.URL) *App {
	a.baseURL = baseURL
	return a
}

// BaseURL returns the domain for the app.
func (a *App) BaseURL() *url.URL {
	return a.baseURL
}

// WithMaxHeaderBytes sets the max header bytes value and returns a reference.
func (a *App) WithMaxHeaderBytes(byteCount int) *App {
	a.maxHeaderBytes = byteCount
	return a
}

// MaxHeaderBytes returns the app max header bytes.
func (a *App) MaxHeaderBytes() int {
	return a.maxHeaderBytes
}

// WithReadHeaderTimeout returns the read header timeout for the server.
func (a *App) WithReadHeaderTimeout(timeout time.Duration) *App {
	a.readHeaderTimeout = timeout
	return a
}

// ReadHeaderTimeout returns the read header timeout for the server.
func (a *App) ReadHeaderTimeout() time.Duration {
	return a.readHeaderTimeout
}

// WithReadTimeout sets the read timeout for the server and returns a reference to the app for building apps with a fluent api.
func (a *App) WithReadTimeout(timeout time.Duration) *App {
	a.readTimeout = timeout
	return a
}

// ReadTimeout returns the read timeout for the server.
func (a *App) ReadTimeout() time.Duration {
	return a.readTimeout
}

// WithIdleTimeout sets the idle timeout.
func (a *App) WithIdleTimeout(timeout time.Duration) *App {
	a.idleTimeout = timeout
	return a
}

// IdleTimeout is the time before we close a connection.
func (a *App) IdleTimeout() time.Duration {
	return a.idleTimeout
}

// WithWriteTimeout sets the write timeout for the server and returns a reference to the app for building apps with a fluent api.
func (a *App) WithWriteTimeout(timeout time.Duration) *App {
	a.writeTimeout = timeout
	return a
}

// WriteTimeout returns the write timeout for the server.
func (a *App) WriteTimeout() time.Duration {
	return a.writeTimeout
}

// WithHSTS enables or disables issuing the strict transport security header.
func (a *App) WithHSTS(hsts *HSTSConfig) *App {
	a.hsts = hsts
	return a
}

// HSTS returns the hsts config.
func (a *App) HSTS() *HSTSConfig {
	return a.hsts
}

// WithTLSConfig sets the tls config for the app.
func (a *App) WithTLSConfig(config *tls.Config) *App {
	a.tls = config
	return a
}

// TLSConfig returns the app tls config.
func (a *App) TLSConfig() *tls.Config {
	return a.tls
}

// SetTLSClientCertPool set the client cert pool from a given set of pems.
func (a *App) SetTLSClientCertPool(certs ...[]byte) error {
	if a.tls == nil {
		a.tls = &tls.Config{}
	}
	a.tls.ClientCAs = x509.NewCertPool()
	for _, cert := range certs {
		ok := a.tls.ClientCAs.AppendCertsFromPEM(cert)
		if !ok {
			return exception.New("invalid ca cert for client cert pool")
		}
	}
	a.tls.BuildNameToCertificate()

	// this forces the server to reload the tls config for every request if there is a cert pool loaded.
	// normally this would introduce overhead but it allows us to hot patch the cert pool.
	a.tls.GetConfigForClient = func(_ *tls.ClientHelloInfo) (*tls.Config, error) {
		return a.tls, nil
	}
	return nil
}

// WithTLSClientCertVerification sets the verification level for client certs.
func (a *App) WithTLSClientCertVerification(verification tls.ClientAuthType) *App {
	if a.tls == nil {
		a.tls = &tls.Config{}
	}
	a.tls.ClientAuth = verification
	return a
}

// WithPort sets the port for the bind address of the app, and returns a reference to the app.
func (a *App) WithPort(port int32) *App {
	a.bindAddr = fmt.Sprintf(":%v", port)
	return a
}

// WithBindAddr sets the address the app listens on, and returns a reference to the app.
func (a *App) WithBindAddr(bindAddr string) *App {
	a.bindAddr = bindAddr
	return a
}

// BindAddr returns the address the server will bind to.
func (a *App) BindAddr() string {
	return a.bindAddr
}

// WithLogger sets the app logger agent and returns a reference to the app.
// It also sets underlying loggers in any child resources like providers and the auth manager.
func (a *App) WithLogger(log *logger.Logger) *App {
	a.log = log
	return a
}

// Logger returns the diagnostics agent for the app.
func (a *App) Logger() *logger.Logger {
	return a.log
}

// WithDefaultMiddlewares sets the application wide default middleware.
func (a *App) WithDefaultMiddlewares(middleware ...Middleware) *App {
	a.defaultMiddleware = middleware
	return a
}

// WithDefaultMiddleware sets the application wide default middleware.
func (a *App) WithDefaultMiddleware(middleware ...Middleware) *App {
	a.defaultMiddleware = append(a.defaultMiddleware, middleware...)
	return a
}

// DefaultMiddleware returns the default middleware.
func (a *App) DefaultMiddleware() []Middleware {
	return a.defaultMiddleware
}

// WithTracer sets the tracer.
func (a *App) WithTracer(tracer Tracer) *App {
	a.tracer = tracer
	return a
}

// Tracer returns the tracer.
func (a *App) Tracer() Tracer {
	return a.tracer
}

// CreateServer returns the basic http.Server for the app.
func (a *App) CreateServer() *http.Server {
	return &http.Server{
		Addr:              a.BindAddr(),
		Handler:           a.Handler(),
		MaxHeaderBytes:    a.maxHeaderBytes,
		ReadTimeout:       a.readTimeout,
		ReadHeaderTimeout: a.readHeaderTimeout,
		WriteTimeout:      a.writeTimeout,
		IdleTimeout:       a.idleTimeout,
		TLSConfig:         a.tls,
	}
}

// WithServer sets the server.
func (a *App) WithServer(server *http.Server) *App {
	a.server = server
	return a
}

// Server returns the underyling http server.
func (a *App) Server() *http.Server {
	return a.server
}

// WithHandler sets the handler.
func (a *App) WithHandler(handler http.Handler) *App {
	a.handler = handler
	return a
}

// Handler returns either a custom handler or the app.
func (a *App) Handler() http.Handler {
	if a.handler != nil {
		return a.handler
	}
	return a
}

// Listener returns the underlying listener.
func (a *App) Listener() *net.TCPListener {
	return a.listener
}

// StartupTasks runs common startup tasks.
func (a *App) StartupTasks() error {
	return a.views.Initialize()
}

// Start starts the server and binds to the given address.
func (a *App) Start() (err error) {
	start := time.Now()
	if a.log != nil {
		a.log.SyncTrigger(NewAppEvent(AppStart).WithApp(a))
		defer a.log.SyncTrigger(NewAppEvent(AppExit).WithApp(a).WithErr(err))
	}

	if a.tls == nil && a.cfg != nil {
		a.tls, err = a.cfg.TLS.GetConfig()
		if err != nil {
			return
		}
	}

	if a.server == nil {
		a.server = a.CreateServer()
	}

	// initialize the view cache.
	err = a.StartupTasks()
	if err != nil {
		return
	}

	serverProtocol := "http"
	if a.server.TLSConfig != nil {
		serverProtocol = "https (tls)"
	}

	a.syncInfof("%s server started, listening on %s", serverProtocol, a.bindAddr)
	if a.log != nil {
		if a.log.Flags() != nil {
			a.syncInfof("%s server logging flags %s", serverProtocol, a.log.Flags().String())
		}
	}

	if a.server.TLSConfig != nil && a.server.TLSConfig.ClientCAs != nil {
		a.syncInfof("%s using client cert pool with (%d) client certs", serverProtocol, len(a.server.TLSConfig.ClientCAs.Subjects()))
	}

	var listener net.Listener
	listener, err = net.Listen("tcp", a.bindAddr)
	if err != nil {
		err = exception.New(err)
		return
	}
	a.listener = listener.(*net.TCPListener)

	if a.log != nil {
		a.log.SyncTrigger(NewAppEvent(AppStartComplete).WithApp(a).WithElapsed(time.Since(start)))
	}

	keepAliveListener := TCPKeepAliveListener{a.listener}
	var shutdownErr error
	a.latch.Started()

	if a.server.TLSConfig != nil {
		shutdownErr = a.server.Serve(tls.NewListener(keepAliveListener, a.server.TLSConfig))
	} else {
		shutdownErr = a.server.Serve(keepAliveListener)
	}
	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		err = exception.New(shutdownErr)
	}

	a.latch.Stopped()
	return
}

// Shutdown is an alias to stop, and stops the server.
func (a *App) Shutdown() error {
	return a.Stop()
}

// Stop stops the server.
func (a *App) Stop() error {
	if !a.Latch().IsRunning() {
		return nil
	}
	a.latch.Stopping()

	ctx := context.Background()
	var cancel context.CancelFunc
	if a.shutdownGracePeriod > 0 {
		ctx, cancel = context.WithTimeout(ctx, a.shutdownGracePeriod)
		defer cancel()
	}
	a.syncInfof("server shutting down")
	a.server.SetKeepAlivesEnabled(false)
	if err := a.server.Shutdown(ctx); err != nil {
		return exception.New(err)
	}

	return nil
}

// WithControllers registers given controllers and returns a reference to the app.
func (a *App) WithControllers(controllers ...Controller) *App {
	for _, c := range controllers {
		a.Register(c)
	}
	return a
}

// Register registers a controller with the app's router.
func (a *App) Register(c Controller) {
	c.Register(a)
}

// --------------------------------------------------------------------------------
// Result Providers
// --------------------------------------------------------------------------------

// WithDefaultResultProvider sets the default result provider.
func (a *App) WithDefaultResultProvider(drp ResultProvider) *App {
	a.defaultResultProvider = drp
	return a
}

// DefaultResultProvider returns the app wide default result provider.
func (a *App) DefaultResultProvider() ResultProvider {
	return a.defaultResultProvider
}

// --------------------------------------------------------------------------------
// Auth Manager
// --------------------------------------------------------------------------------

// WithAuth sets the auth manager.
func (a *App) WithAuth(am *AuthManager) *App {
	a.auth = am
	return a
}

// Auth returns the session manager.
func (a *App) Auth() *AuthManager {
	return a.auth
}

// --------------------------------------------------------------------------------
// Views
// --------------------------------------------------------------------------------

// WithViews sets the view cache.
func (a *App) WithViews(vc *ViewCache) *App {
	a.views = vc
	return a
}

// Views returns the view cache.
func (a *App) Views() *ViewCache {
	return a.views
}

// --------------------------------------------------------------------------------
// Static Result Methods
// --------------------------------------------------------------------------------

// SetStaticRewriteRule adds a rewrite rule for a specific statically served path.
// It mutates the path for the incoming static file request to the fileserver according to the action.
func (a *App) SetStaticRewriteRule(route, match string, action RewriteAction) error {
	mountedRoute := a.createStaticMountRoute(route)
	if static, hasRoute := a.statics[mountedRoute]; hasRoute {
		return static.AddRewriteRule(match, action)
	}
	return exception.New("no static fileserver mounted at route").WithMessagef("route: %s", route)
}

// SetStaticHeader adds a header for the given static path.
// These headers are automatically added to any result that the static path fileserver sends.
func (a *App) SetStaticHeader(route, key, value string) error {
	mountedRoute := a.createStaticMountRoute(route)
	if static, hasRoute := a.statics[mountedRoute]; hasRoute {
		static.AddHeader(key, value)
		return nil
	}
	return exception.New("no static fileserver mounted at route").WithMessagef("route: %s", mountedRoute)
}

// SetStaticMiddleware adds static middleware for a given route.
func (a *App) SetStaticMiddleware(route string, middlewares ...Middleware) error {
	mountedRoute := a.createStaticMountRoute(route)
	if static, hasRoute := a.statics[mountedRoute]; hasRoute {
		static.SetMiddleware(middlewares...)
		return nil
	}
	return exception.New("no static fileserver mounted at route").WithMessagef("route: %s", mountedRoute)
}

// ServeStatic serves files from the given file system root.
// If the path does not end with "/*filepath" that suffix will be added for you internally.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
func (a *App) ServeStatic(route, filepath string) {
	sfs := NewStaticFileServer(http.Dir(filepath))
	mountedRoute := a.createStaticMountRoute(route)
	a.statics[mountedRoute] = sfs
	a.Handle("GET", mountedRoute, a.renderAction(a.middlewarePipeline(sfs.Action)))
}

// ServeStaticCached serves files from the given file system root.
// If the path does not end with "/*filepath" that suffix will be added for you internally.
func (a *App) ServeStaticCached(route, filepath string) {
	sfs := NewCachedStaticFileServer(http.Dir(filepath))
	mountedRoute := a.createStaticMountRoute(route)
	a.statics[mountedRoute] = sfs
	a.Handle("GET", mountedRoute, a.renderAction(a.middlewarePipeline(sfs.Action)))
}

func (a *App) createStaticMountRoute(route string) string {
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

// --------------------------------------------------------------------------------
// Router internal methods
// --------------------------------------------------------------------------------

// WithNotFoundHandler sets the not found handler.
func (a *App) WithNotFoundHandler(handler Action) *App {
	a.notFoundHandler = a.renderAction(handler)
	return a
}

// NotFoundHandler returns the not found handler.
func (a *App) NotFoundHandler() Handler {
	return a.notFoundHandler
}

// WithMethodNotAllowedHandler sets the not allowed handler.
func (a *App) WithMethodNotAllowedHandler(handler Action) *App {
	a.methodNotAllowedHandler = a.renderAction(handler)
	return a
}

// MethodNotAllowedHandler returns the method not allowed handler.
func (a *App) MethodNotAllowedHandler() Handler {
	return a.methodNotAllowedHandler
}

// WithPanicAction sets the panic action.
func (a *App) WithPanicAction(action PanicAction) *App {
	a.panicAction = action
	return a
}

// PanicAction returns the panic action.
func (a *App) PanicAction() PanicAction {
	return a.panicAction
}

// --------------------------------------------------------------------------------
// Testing Methods
// --------------------------------------------------------------------------------

// Mock returns a request bulider to facilitate mocking requests against the app
// without having to start it and bind it to a port.
/*
An example mock request that hits an already registered "GET" route at "/foo":

	assert.Nil(app.Mock().Get("/").Execute())

This will assert that the request completes successfully, but does not return the
response.
*/
func (a *App) Mock() *MockRequestBuilder {
	return NewMockRequestBuilder(a).WithErr(a.StartupTasks())
}

// --------------------------------------------------------------------------------
// Route Registration / HTTP Methods
// --------------------------------------------------------------------------------

// GET registers a GET request handler.
/*
Routes should be registered in the form:

	app.GET("/myroute", myAction, myMiddleware...)

It is important to note that routes are registered in order and
cannot have any wildcards inside the routes.
*/
func (a *App) GET(path string, action Action, middleware ...Middleware) {
	a.Handle("GET", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// OPTIONS registers a OPTIONS request handler.
func (a *App) OPTIONS(path string, action Action, middleware ...Middleware) {
	a.Handle("OPTIONS", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// HEAD registers a HEAD request handler.
func (a *App) HEAD(path string, action Action, middleware ...Middleware) {
	a.Handle("HEAD", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// PUT registers a PUT request handler.
func (a *App) PUT(path string, action Action, middleware ...Middleware) {
	a.Handle("PUT", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// PATCH registers a PATCH request handler.
func (a *App) PATCH(path string, action Action, middleware ...Middleware) {
	a.Handle("PATCH", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// POST registers a POST request actions.
func (a *App) POST(path string, action Action, middleware ...Middleware) {
	a.Handle("POST", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// DELETE registers a DELETE request handler.
func (a *App) DELETE(path string, action Action, middleware ...Middleware) {
	a.Handle("DELETE", path, a.renderAction(a.middlewarePipeline(action, middleware...)))
}

// Handle adds a raw handler at a given method and path.
func (a *App) Handle(method, path string, handler Handler) {
	if len(path) == 0 {
		panic("path must not be empty")
	}
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if a.routes == nil {
		a.routes = make(map[string]*node)
	}

	root := a.routes[method]
	if root == nil {
		root = new(node)
		a.routes[method] = root
	}

	root.addRoute(method, path, handler)
}

// Lookup finds the route data for a given method and path.
func (a *App) Lookup(method, path string) (route *Route, params RouteParameters, slashRedirect bool) {
	if root := a.routes[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// --------------------------------------------------------------------------------
// Request Pipeline
// --------------------------------------------------------------------------------

// ServeHTTP makes the router implement the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if a.recoverPanics {
		defer a.recover(w, req)
	}

	path := req.URL.Path
	if root := a.routes[req.Method]; root != nil {
		if route, params, tsr := root.getValue(path); route != nil {
			route.Handler(w, req, route, params)
			return
		} else if req.Method != MethodConnect && path != "/" {
			code := http.StatusMovedPermanently // 301 // Permanent redirect, request with GET method
			if req.Method != MethodGet {
				code = http.StatusTemporaryRedirect // 307
			}

			if tsr && a.redirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}
		}
	}

	if req.Method == MethodOptions {
		// Handle OPTIONS requests
		if a.handleOptions {
			if allow := a.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set(HeaderAllow, allow)
				return
			}
		}
	} else {
		// Handle 405
		if a.handleMethodNotAllowed {
			if allow := a.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set(HeaderAllow, allow)
				if a.methodNotAllowedHandler != nil {
					a.methodNotAllowedHandler(w, req, nil, nil)
				} else {
					http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				}
				return
			}
		}
	}

	// Handle 404
	if a.notFoundHandler != nil {
		a.notFoundHandler(w, req, nil, nil)
	} else {
		http.NotFound(w, req)
	}
}

// renderAction is the translation step from Action to Handler.
// this is where the bulk of the "pipeline" happens.
func (a *App) renderAction(action Action) Handler {
	return func(w http.ResponseWriter, r *http.Request, route *Route, p RouteParameters) {
		var err error
		var tf TraceFinisher

		var response ResponseWriter
		if strings.Contains(r.Header.Get(HeaderAcceptEncoding), ContentEncodingGZIP) {
			w.Header().Set(HeaderContentEncoding, ContentEncodingGZIP)
			response = NewCompressedResponseWriter(w)
		} else {
			w.Header().Set(HeaderContentEncoding, ContentEncodingIdentity)
			response = NewRawResponseWriter(w)
		}

		ctx := a.createCtx(response, r, route, p)
		ctx.onRequestStart()
		if a.tracer != nil {
			tf = a.tracer.Start(ctx)
		}
		if a.log != nil {
			a.log.Trigger(a.httpRequestEvent(ctx))
		}

		if len(a.defaultHeaders) > 0 {
			for key, value := range a.defaultHeaders {
				response.Header().Set(key, value)
			}
		}

		if a.hsts.GetEnabled() {
			a.addHSTSHeader(response)
		}

		result := action(ctx)
		if result != nil {

			// check for a prerender step
			if typed, ok := result.(ResultPreRender); ok {
				if preRender := typed.PreRender(ctx); preRender != nil {
					err = exception.Nest(err, preRender)
					a.logFatal(err, r)
				}
			}

			// do the render
			a.logError(result.Render(ctx))

			// check for a render complete step
			if typed, ok := result.(ResultRenderComplete); ok {
				if renderComplete := typed.RenderComplete(ctx); renderComplete != nil {
					err = exception.Nest(err, renderComplete)
					a.logFatal(renderComplete, r)
				}
			}
		}

		ctx.onRequestFinish()
		a.logError(response.Close())

		// effectively "request complete"
		if a.log != nil {
			a.log.Trigger(a.httpResponseEvent(ctx))
		}
		if tf != nil {
			tf.Finish(ctx, err)
		}
	}
}

func (a *App) createCtx(w ResponseWriter, r *http.Request, route *Route, p RouteParameters) *Ctx {
	return NewCtx(w, r).
		WithApp(a).
		WithRoute(route).
		WithRouteParams(p).
		WithState(a.state.Copy()).
		WithTracer(a.tracer).
		WithViews(a.views).
		WithAuth(a.auth).
		WithLogger(a.log).
		WithDefaultResultProvider(a.defaultResultProvider)
}

func (a *App) middlewarePipeline(action Action, middleware ...Middleware) Action {
	if len(middleware) == 0 && len(a.defaultMiddleware) == 0 {
		return action
	}

	finalMiddleware := make([]Middleware, len(middleware)+len(a.defaultMiddleware))
	cursor := len(finalMiddleware) - 1
	for i := len(a.defaultMiddleware) - 1; i >= 0; i-- {
		finalMiddleware[cursor] = a.defaultMiddleware[i]
		cursor--
	}

	for i := len(middleware) - 1; i >= 0; i-- {
		finalMiddleware[cursor] = middleware[i]
		cursor--
	}

	return NestMiddleware(action, finalMiddleware...)
}

func (a *App) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range a.routes {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
		return
	}
	for method := range a.routes {
		// Skip the requested method - we already tried this one
		if method == reqMethod || method == "OPTIONS" {
			continue
		}

		handle, _, _ := a.routes[method].getValue(path)
		if handle != nil {
			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

func (a *App) addHSTSHeader(w http.ResponseWriter) {
	parts := []string{fmt.Sprintf(HSTSMaxAgeFormat, a.hsts.GetMaxAgeSeconds())}
	if a.hsts.GetIncludeSubDomains() {
		parts = append(parts, HSTSIncludeSubDomains)
	}
	if a.hsts.GetPreload() {
		parts = append(parts, HSTSPreload)
	}
	w.Header().Set(HeaderStrictTransportSecurity, strings.Join(parts, "; "))
}

func (a *App) httpRequestEvent(ctx *Ctx) *logger.HTTPRequestEvent {
	event := logger.NewHTTPRequestEvent(ctx.Request())
	event.SetEntity(ctx.ID())
	if ctx.Route() != nil {
		event = event.WithRoute(ctx.Route().String())
	}
	return event
}

func (a *App) httpResponseEvent(ctx *Ctx) *logger.HTTPResponseEvent {
	event := logger.NewHTTPResponseEvent(ctx.Request()).
		WithStatusCode(ctx.Response().StatusCode()).
		WithElapsed(ctx.Elapsed()).
		WithContentLength(ctx.Response().ContentLength())
	event.SetEntity(ctx.ID())

	if ctx.Route() != nil {
		event = event.WithRoute(ctx.Route().String())
	}

	if ctx.Response().Header() != nil {
		event = event.WithContentType(ctx.Response().Header().Get(HeaderContentType))
		event = event.WithContentEncoding(ctx.Response().Header().Get(HeaderContentEncoding))
	}
	return event
}

func (a *App) recover(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		err := exception.New(rcv)
		a.logFatal(err, req)
		if a.panicAction != nil {
			a.handlePanic(w, req, rcv)
		} else {
			http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
		}
	}
}

func (a *App) handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
	a.renderAction(func(ctx *Ctx) Result {
		if a.log != nil {
			a.log.Fatalf("%v", err)
		}
		return a.panicAction(ctx, err)
	})(w, r, nil, nil)
}

func (a *App) logFatal(err error, req *http.Request) {
	if a.log == nil {
		return
	}
	if err != nil {
		a.log.FatalWithReq(err, req)
	}
}

func (a *App) logError(err error) {
	if a.log == nil {
		return
	}
	if err != nil {
		a.log.Error(err)
	}
}

func (a *App) syncInfof(format string, args ...interface{}) {
	if a.log == nil {
		return
	}
	a.log.SyncInfof(format, args...)
}

func (a *App) syncFatalf(format string, args ...interface{}) {
	if a.log == nil {
		return
	}
	a.log.SyncFatalf(format, args...)
}
