package main

import (
	"os"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All()
	app := web.New(web.OptLog(log))
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
		log.Fatal(err)
		os.Exit(1)
	}
}
