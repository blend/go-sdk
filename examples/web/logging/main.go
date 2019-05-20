package main

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/graceful"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All()
	app := web.New(web.OptLog(log))
	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Text.Result("foo")
	})
	log.Listen(logger.HTTPRequest, logger.DefaultListenerName, logger.NewHTTPRequestEventListener(func(_ context.Context, wre *logger.HTTPRequestEvent) {
		fmt.Printf("Route: %s\n", wre.Route)
	}))

	graceful.Shutdown(app)
}
