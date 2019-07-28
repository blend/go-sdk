package main

import (
	"fmt"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	app := web.MustNew(web.OptLog(logger.All()))

	app.GET("/204", func(_ *web.Ctx) web.Result {
		return web.NoContent
	})
	app.GET("/500", func(_ *web.Ctx) web.Result {
		return web.JSON.InternalError(fmt.Errorf("this is only a test"))
	})

	if err := graceful.Shutdown(app); err != nil {
		logger.FatalExit(err)
	}
}
