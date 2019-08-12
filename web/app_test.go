package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
)

// assert an app is graceful
var (
	_ graceful.Graceful = (*App)(nil)
)

func controllerNoOp(_ *Ctx) Result { return nil }

type testController struct {
	callback func(app *App)
}

func (tc testController) Register(app *App) {
	if tc.callback != nil {
		tc.callback(app)
	}
}

func TestAppNew(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)
	assert.NotNil(app.State)
	assert.NotNil(app.Views)
}

func TestAppNewFromConfig(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptConfig(Config{
		BindAddr:               ":5555",
		Port:                   5000,
		HandleMethodNotAllowed: true,
		HandleOptions:          true,
		DisablePanicRecovery:   true,

		MaxHeaderBytes:    128,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       6 * time.Second,
		IdleTimeout:       7 * time.Second,
		WriteTimeout:      8 * time.Second,

		CookieName: "A GOOD ONE",
		Views: ViewCacheConfig{
			LiveReload: true,
		},
	}))

	assert.Nil(err)

	assert.Equal(":5555", app.Config.BindAddr)
	assert.True(app.Config.HandleMethodNotAllowed)
	assert.True(app.Config.HandleOptions)
	assert.True(app.Config.DisablePanicRecovery)
	assert.Equal(128, app.Config.MaxHeaderBytes)
	assert.Equal(5*time.Second, app.Config.ReadHeaderTimeout)
	assert.Equal(6*time.Second, app.Config.ReadTimeout)
	assert.Equal(7*time.Second, app.Config.IdleTimeout)
	assert.Equal(8*time.Second, app.Config.WriteTimeout)
	assert.Equal("A GOOD ONE", app.Auth.CookieDefaults.Name, "we should use the auth config for the auth manager")
	assert.True(app.Views.LiveReload, "we should use the view cache config for the view cache")
}

func TestAppRegister(t *testing.T) {
	assert := assert.New(t)
	called := false
	c := &testController{
		callback: func(_ *App) {
			called = true
		},
	}
	app, err := New()

	assert.Nil(err)
	assert.False(called)
	app.Register(c)
	assert.True(called)
}

func TestAppPathParams(t *testing.T) {
	assert := assert.New(t)

	var route *Route
	var params RouteParameters
	app, err := New()
	assert.Nil(err)
	app.GET("/:uuid", func(c *Ctx) Result {
		route = c.Route
		params = c.RouteParams
		return Raw([]byte("ok!"))
	})

	route, params, skipSlashRedirect := app.Lookup("GET", "/foo")
	assert.NotNil(route)
	assert.NotEmpty(params)
	assert.Equal("foo", params.Get("uuid"))
	assert.False(skipSlashRedirect)

	meta, err := MockGet(app, "/foo").DiscardWithResponse()
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotNil(route)
	assert.Equal("GET", route.Method)
	assert.Equal("/:uuid", route.Path)
	assert.NotNil(route.Handler)

	assert.NotEmpty(params)
	assert.Equal("foo", params.Get("uuid"))
}

func TestAppPathParamsForked(t *testing.T) {
	/*
		this test should assert that we can have a common structure of routes
		namely that you can have some shared prefix but differentiate by plural.
	*/

	assert := assert.New(t)

	var route *Route
	var params RouteParameters
	app, err := New()
	assert.Nil(err)
	app.GET("/foo/:uuid", func(c *Ctx) Result { return NoContent })
	app.GET("/foos/bar/:uuid", func(c *Ctx) Result {
		route = c.Route
		params = c.RouteParams
		return Raw([]byte("ok!"))
	})

	meta, err := MockGet(app, "/foos/bar/foo").DiscardWithResponse()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotNil(route)
	assert.Equal("GET", route.Method)
	assert.Equal("/foos/bar/:uuid", route.Path)
	assert.NotNil(route.Handler)

	assert.NotNil(params)
	assert.NotEmpty(params)
	assert.Equal("foo", params.Get("uuid"))
}

func TestAppSetLogger(t *testing.T) {
	assert := assert.New(t)

	log := logger.MustNew()
	app, err := New(OptLog(log))
	assert.Nil(err)
	assert.NotNil(app.Log)
}

func TestAppCreateStaticMountedRoute(t *testing.T) {
	assert := assert.New(t)
	app, err := New()
	assert.Nil(err)
	assert.Equal("/testPath/*filepath", app.formatStaticMountRoute("/testPath/*filepath"))
	assert.Equal("/testPath/*filepath", app.formatStaticMountRoute("/testPath/"))
	assert.Equal("/testPath/*filepath", app.formatStaticMountRoute("/testPath"))
}

func TestAppStaticRewrite(t *testing.T) {
	assert := assert.New(t)
	app, err := New()
	assert.Nil(err)

	app.ServeStatic("/testPath", []string{"_static"})
	assert.NotEmpty(app.Statics)
	assert.NotNil(app.Statics["/testPath/*filepath"])
	assert.Nil(app.SetStaticRewriteRule("/testPath", "(.*)", func(path string, pieces ...string) string {
		return path
	}))
	assert.NotNil(app.SetStaticRewriteRule("/notapath", "(.*)", func(path string, pieces ...string) string {
		return path
	}))

	assert.NotEmpty(app.Statics["/testPath/*filepath"].RewriteRules)
}

func TestAppStaticRewriteBadExp(t *testing.T) {
	assert := assert.New(t)
	app, err := New()
	assert.Nil(err)

	app.ServeStatic("/testPath", []string{"_static"})
	assert.NotEmpty(app.Statics)
	assert.NotNil(app.Statics["/testPath/*filepath"])

	err = app.SetStaticRewriteRule("/testPath", "((((", func(path string, pieces ...string) string {
		return path
	})

	assert.NotNil(err)
	assert.Empty(app.Statics["/testPath/*filepath"].RewriteRules)
}

func TestAppStaticHeader(t *testing.T) {
	assert := assert.New(t)
	app, err := New()
	assert.Nil(err)

	app.ServeStatic("/testPath", []string{"_static"})
	assert.NotEmpty(app.Statics)
	assert.NotNil(app.Statics["/testPath/*filepath"])
	assert.Nil(app.SetStaticHeader("/testPath/*filepath", "cache-control", "haha what is caching."))
	assert.NotNil(app.SetStaticHeader("/notaroute", "cache-control", "haha what is caching."))
	assert.NotEmpty(app.Statics["/testPath/*filepath"].Headers)
}

func TestAppMiddleWarePipeline(t *testing.T) {
	assert := assert.New(t)

	didRun := false

	app, err := New()
	assert.Nil(err)

	app.GET("/",
		func(r *Ctx) Result { return Raw([]byte("OK!")) },
		func(action Action) Action {
			didRun = true
			return action
		},
		func(action Action) Action {
			return func(r *Ctx) Result {
				return Raw([]byte("foo"))
			}
		},
	)

	result, err := MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.True(didRun)
	assert.Equal("foo", string(result))
}

func TestAppStatic(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)

	app.ServeStatic("/static/*filepath", []string{"testdata"})

	index, err := MockGet(app, "/static/test_file.html").Bytes()
	assert.Nil(err)
	assert.True(strings.Contains(string(index), "Test!"), string(index))
}

func TestAppStaticSingleFile(t *testing.T) {
	assert := assert.New(t)
	app, err := New()
	assert.Nil(err)

	app.GET("/", func(r *Ctx) Result {
		return Static("testdata/test_file.html")
	})

	index, err := MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.True(strings.Contains(string(index), "Test!"), string(index))
}

func TestAppProviderMiddleware(t *testing.T) {
	assert := assert.New(t)

	var okAction = func(r *Ctx) Result {
		assert.NotNil(r.DefaultProvider)
		return Raw([]byte("OK!"))
	}

	app, err := New()
	assert.Nil(err)

	app.GET("/", okAction, JSONProviderAsDefault)

	err = MockGet(app, "/").Discard()
	assert.Nil(err)
}

func TestAppProviderMiddlewareOrder(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)

	var okAction = func(r *Ctx) Result {
		assert.NotNil(r.DefaultProvider)
		return Raw([]byte("OK!"))
	}

	var dependsOnProvider = func(action Action) Action {
		return func(r *Ctx) Result {
			assert.NotNil(r.DefaultProvider)
			return action(r)
		}
	}

	app.GET("/", okAction, dependsOnProvider, JSONProviderAsDefault)
	assert.Nil(MockGet(app, "/").Discard())
}

func TestAppDefaultResultProvider(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptUse(ViewProviderAsDefault))
	assert.Nil(err)

	assert.NotEmpty(app.DefaultMiddleware)
	rc := app.createCtx(nil, nil, nil, nil)
	assert.NotNil(rc.DefaultProvider)
	assert.NotNil(rc.App)
}

func TestAppDefaultResultProviderWithDefault(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptUse(ViewProviderAsDefault))
	assert.Nil(err)
	assert.NotEmpty(app.DefaultMiddleware)

	rc := app.createCtx(nil, nil, nil, nil)

	// this will be set to the default initially
	assert.NotNil(rc.DefaultProvider)

	app.GET("/", func(ctx *Ctx) Result {
		assert.NotNil(ctx.DefaultProvider)
		_, isTyped := ctx.DefaultProvider.(*ViewCache)
		assert.True(isTyped)
		return nil
	})
	assert.Nil(MockGet(app, "/").Discard())
}

func TestAppDefaultResultProviderWithDefaultFromRoute(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptUse(JSONProviderAsDefault))
	assert.Nil(err)

	app.Views.AddLiterals(DefaultTemplateNotAuthorized)
	app.GET("/", controllerNoOp, SessionRequired, ViewProviderAsDefault)

	//somehow assert that the content type is html
	meta, err := MockGet(app, "/").DiscardWithResponse()
	assert.Nil(err)
	defer meta.Body.Close()

	assert.Equal(ContentTypeHTML, meta.Header.Get(HeaderContentType))
}

func TestAppViewResult(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)

	app.Views.AddPaths("testdata/test_file.html")
	app.GET("/", func(r *Ctx) Result {
		return r.Views.View("test", "foobarbaz")
	})

	contents, meta, err := MockGet(app, "/").BytesWithResponse()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Equal(ContentTypeHTML, meta.Header.Get(HeaderContentType))
	assert.Contains(string(contents), "foobarbaz")
}

func TestAppWritesLogs(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	agent := logger.MustNew(logger.OptAll(), logger.OptOutput(buffer))

	app, err := New(OptLog(agent))
	assert.Nil(err)

	app.GET("/", func(r *Ctx) Result {
		return Raw([]byte("ok!"))
	})
	err = MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Nil(agent.Drain())

	assert.NotZero(buffer.Len())
	assert.NotEmpty(buffer.String())
}

func TestAppBindAddr(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvironmentVariableBindAddr, ":9999")
	env.Env().Set(EnvironmentVariablePort, "1111")
	defer env.Restore()

	assert.Equal(":3333", MustNew(OptBindAddr(":3333")).Config.BindAddr)
	assert.Equal(":2222", MustNew(OptPort(2222)).Config.BindAddr)
}

func TestAppNotFound(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)

	app.GET("/", func(r *Ctx) Result {
		return Raw([]byte("ok!"))
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	app.NotFoundHandler = app.RenderAction(func(r *Ctx) Result {
		defer wg.Done()
		return JSON.NotFound()
	})
	err = MockGet(app, "/doesntexist").Discard()
	assert.Nil(err)
	wg.Wait()
}

func TestAppDefaultHeaders(t *testing.T) {
	assert := assert.New(t)
	app, err := New(OptDefaultHeader("foo", "bar"), OptDefaultHeader("baz", "buzz"))
	assert.Nil(err)
	assert.Equal([]string{"buzz"}, app.DefaultHeaders[http.CanonicalHeaderKey("baz")])

	app.GET("/", func(r *Ctx) Result {
		return Text.Result("ok")
	})

	meta, err := MockGet(app, "/").DiscardWithResponse()
	assert.Nil(err)
	assert.NotEmpty(meta.Header)
	assert.Equal("bar", meta.Header.Get("foo"))
	assert.Equal("buzz", meta.Header.Get("baz"))
}

func TestAppViewErrorsRenderErrorView(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)
	assert.NotNil(app.Views)

	app.Views.AddLiterals(`{{ define "malformed" }} {{ .Ctx ALSKADJALSKDJA }} {{ end }}`)
	app.GET("/", func(r *Ctx) Result {
		return r.Views.View("malformed", nil)
	})
	assert.NotNil(MockGet(app, "/").Discard())
}

func TestAppAddsDefaultHeaders(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptBindAddr(DefaultMockBindAddr))
	assert.Nil(err)

	app.GET("/", func(r *Ctx) Result {
		return Text.Result("OK!")
	})

	go app.Start()
	<-app.NotifyStarted()
	defer app.Stop()

	res, err := http.Get("http://" + app.Listener.Addr().String() + "/")
	assert.Nil(err)
	assert.NotEmpty(res.Header)
	assert.Equal(PackageName, res.Header.Get(HeaderServer))
}

func TestAppHandlesPanics(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptBindAddr(DefaultMockBindAddr))
	assert.Nil(err)

	app.GET("/", func(r *Ctx) Result {
		panic("this is only a test")
	})

	var didRecover bool
	go func() {
		defer func() {
			if r := recover(); r != nil {
				didRecover = true
			}
		}()
		app.Start()
	}()
	defer app.Stop()
	<-app.NotifyStarted()

	res, err := http.Get("http://" + app.Listener.Addr().String() + "/")
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, res.StatusCode)
	assert.False(didRecover)
}

var (
	_ Tracer     = (*mockTracer)(nil)
	_ ViewTracer = (*mockTracer)(nil)
)

type mockTracer struct {
	OnStart  func(*Ctx)
	OnFinish func(*Ctx, error)

	OnViewStart  func(*Ctx, *ViewResult)
	OnViewFinish func(*Ctx, *ViewResult, error)
}

func (mt mockTracer) Start(ctx *Ctx) TraceFinisher {
	if mt.OnStart != nil {
		mt.OnStart(ctx)
	}
	return &mockTraceFinisher{parent: &mt}
}

func (mt mockTracer) StartView(ctx *Ctx, vr *ViewResult) ViewTraceFinisher {
	if mt.OnViewStart != nil {
		mt.OnViewStart(ctx, vr)
	}
	return &mockViewTraceFinisher{parent: &mt}
}

type mockTraceFinisher struct {
	parent *mockTracer
}

func (mtf mockTraceFinisher) Finish(ctx *Ctx, err error) {
	mtf.parent.OnFinish(ctx, err)
}

type mockViewTraceFinisher struct {
	parent *mockTracer
}

func (mvf mockViewTraceFinisher) FinishView(ctx *Ctx, vr *ViewResult, err error) {
	mvf.parent.OnViewFinish(ctx, vr, err)
}

func ok(_ *Ctx) Result            { return JSON.OK() }
func doPanic(_ *Ctx) Result       { panic("this is only a test") }
func internalError(_ *Ctx) Result { return JSON.InternalError(fmt.Errorf("only a test")) }
func viewOK(ctx *Ctx) Result      { return ctx.Views.View("ok", nil) }

func TestAppTracer(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var hasValue bool

	app, err := New()
	assert.Nil(err)

	app.GET("/", ok)
	app.Tracer = mockTracer{
		OnStart: func(ctx *Ctx) {
			defer wg.Done()
			ctx.WithStateValue("foo", "bar")
		},
		OnFinish: func(ctx *Ctx, err error) {
			defer wg.Done()
			hasValue = ctx.StateValue("foo") != nil
		},
	}

	assert.Nil(MockGet(app, "/").Discard())
	wg.Wait()

	assert.True(hasValue)
}

func TestAppTracerError(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)

	var hasError bool

	app, err := New()
	assert.Nil(err)

	app.GET("/", ok)
	app.GET("/error", internalError)
	app.Tracer = mockTracer{
		OnFinish: func(ctx *Ctx, err error) {
			defer wg.Done()
			hasError = err != nil
		},
	}

	assert.Nil(MockGet(app, "/error").Discard())
	wg.Wait()
	assert.True(hasError)
}

func TestAppViewTracer(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)

	var hasValue bool

	app, err := New()
	assert.Nil(err)

	app.Views.AddLiterals("{{ define \"ok\" }}ok{{end}}")
	assert.Nil(app.Views.Initialize())

	app.GET("/", ok)
	app.GET("/view", viewOK)
	app.Tracer = mockTracer{
		OnStart:  func(_ *Ctx) { wg.Done() },
		OnFinish: func(_ *Ctx, _ error) { wg.Done() },
		OnViewStart: func(ctx *Ctx, vr *ViewResult) {
			defer wg.Done()
			hasValue = vr.ViewName == "ok"
		},
		OnViewFinish: func(ctx *Ctx, vr *ViewResult, err error) {
			defer wg.Done()
		},
	}

	assert.Nil(MockGet(app, "/view").Discard())
	wg.Wait()

	assert.True(hasValue)
}

func TestAppViewTracerError(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)

	var hasValue, hasError, hasViewError bool

	app, err := New()
	assert.Nil(err)

	app.Views.AddLiterals("{{ define \"ok\" }}{{template \"fake\"}}ok{{end}}")
	app.GET("/view", viewOK)
	app.Tracer = mockTracer{
		OnStart: func(_ *Ctx) { wg.Done() },
		OnFinish: func(_ *Ctx, err error) {
			defer wg.Done()
			hasError = err != nil
		},
		OnViewStart: func(ctx *Ctx, vr *ViewResult) {
			defer wg.Done()
			hasValue = vr.ViewName == "ok"
		},
		OnViewFinish: func(ctx *Ctx, vr *ViewResult, err error) {
			defer wg.Done()
			hasViewError = err != nil
		},
	}

	assert.Nil(MockGet(app, "/view").Discard())
	wg.Wait()

	assert.True(hasValue)
	assert.True(hasError)
	assert.True(hasViewError)
}

func TestAppAllowed(t *testing.T) {
	assert := assert.New(t)

	app, err := New()
	assert.Nil(err)
	app.GET("/test", nil)

	allowed := strings.Split(app.allowed("*", ""), ", ")
	assert.Len(allowed, 1)
	assert.Equal("GET", allowed[0])

	app.POST("/hello", nil)
	allowed = strings.Split(app.allowed("*", ""), ", ")
	assert.Len(allowed, 2)
	assert.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == "GET"
	})
	assert.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == "POST"
	})

	app, err = New()
	assert.Nil(err)

	app.GET("/hello", controllerNoOp)
	allowed = strings.Split(app.allowed("/hello", ""), ", ")
	assert.Len(allowed, 2)
	assert.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == "GET"
	})
	assert.Any(allowed, func(i interface{}) bool {
		s, ok := i.(string)
		return ok && s == "OPTIONS"
	})
	app.POST("/hello", controllerNoOp)
	allowed = strings.Split(app.allowed("/hello", ""), ", ")
	assert.Len(allowed, 3)

	app.OPTIONS("/hello", controllerNoOp)
	app.HEAD("/hello", controllerNoOp)
	app.PUT("/hello", controllerNoOp)
	app.DELETE("/hello", controllerNoOp)

	app.PATCH("/hi", controllerNoOp)
	app.PATCH("/there", controllerNoOp)
	allowed = strings.Split(app.allowed("/hello", ""), ", ")
	assert.Len(allowed, 6)

	app.PATCH("/hello", controllerNoOp)
	allowed = strings.Split(app.allowed("/hello", ""), ", ")
	assert.Len(allowed, 7)
}

func TestAppNilLoggerPanic(t *testing.T) {
	assert := assert.New(t)

	app, err := New(OptLog(nil))
	assert.Nil(err)
	app.PanicAction = func(r *Ctx, err interface{}) Result {
		return r.DefaultProvider.InternalError(ex.New(err))
	}
	assert.Nil(err)
	app.GET("/", doPanic, ViewProviderAsDefault)

	res, err := MockGet(app, "/").DiscardWithResponse()
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, res.StatusCode)
}
