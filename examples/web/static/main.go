package main

import (
	"net/http"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All()
	app := web.MustNew(web.OptLog(log))
	csf := web.NewStaticFileServer(
		web.OptStaticFileServerSearchPaths(http.Dir(".")),
	)

	app.ServeStatic("/static/*filepath", []string{"_static"})
	app.ServeStaticCached("/static_cached/*filepath", []string{"_static"})
	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Static("index.html")
	})
	app.GET("/cached", func(r *web.Ctx) web.Result {
		return csf.ServeFile(r, "index.html")
	})
	graceful.Shutdown(app)
}
