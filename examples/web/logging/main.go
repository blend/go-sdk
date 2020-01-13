package main

import (
	"context"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/webutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.Prod()
	app := web.MustNew(web.OptLog(log))
	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Text.Result("foo")
	})
	log.Listen(webutil.HTTPRequest, logger.DefaultListenerName, webutil.NewHTTPRequestEventListener(func(_ context.Context, wre webutil.HTTPRequestEvent) {
		log.Infof("got a new request at route: %s", wre.Route)
	}))

	graceful.Shutdown(app)
}
