package main

import (
	"os"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	app := web.MustNew(web.OptLog(logger.All(logger.OptHeading("views-test"))))
	app.Views.AddPaths(
		"_views/header.html",
		"_views/footer.html",
		"_views/index.html",
	)

	app.Views.FuncMap["foo"] = func() string {
		return "hello!"
	}

	if len(os.Getenv("LIVE_RELOAD")) > 0 {
		app.Views.LiveReload = true
	}

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Views.View("index", nil)
	})

	if err := graceful.Shutdown(app); err != nil {
		app.Log.Fatal(err)
		os.Exit(1)
	}
}
