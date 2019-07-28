package main

import (
	"fmt"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	app := web.MustNew(web.OptBindAddr(":8080"), web.OptLog(logger.All()))
	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Text.Result("ok!")
	})

	app.POST("/reparse", func(r *web.Ctx) web.Result {
		body, err := r.PostBody()
		if err != nil {
			return web.Text.BadRequest(err)
		}
		if len(body) == 0 {
			return web.Text.BadRequest(fmt.Errorf("empty body"))
		}
		return web.Text.Result(web.StringValue(r.Param("foo")))
	})
	if err := graceful.Shutdown(app); err != nil {
		logger.FatalExit(err)
	}
}
