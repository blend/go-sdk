package main

import (
	"os"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.Prod()
	app := web.MustNew(
		web.OptLog(log),
		web.OptConfigFromEnv(),
		web.OptUse(web.GZip),
	)
	app.GET("/", func(_ *web.Ctx) web.Result { return web.Text.Result("OK!") })
	if err := graceful.Shutdown(app); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
