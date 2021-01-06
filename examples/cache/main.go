package main

import (
	"os"
	"time"

	"github.com/blend/go-sdk/cache"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"
)

func getData() (interface{}, error) {
	time.Sleep(500 * time.Millisecond)
	var output []string
	for x := 0; x < 1024; x++ {
		output = append(output, uuid.V4().String())
	}
	return output, nil
}

func main() {
	log := logger.Prod()
	log.Disable(webutil.FlagHTTPRequest) // disable noisey events.
	app, err := web.New(
		web.OptConfigFromEnv(),
		web.OptLog(log),
		web.OptUse(web.GZip), // NOTE: as of v3.0.0 gzip response compression middleware is not enabled by default, you _must_ enable it explicitly.
		web.OptShutdownGracePeriod(time.Second),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	app.PanicAction = func(_ *web.Ctx, r interface{}) web.Result {
		return web.Text.InternalError(ex.New(r))
	}

	lc := cache.New(cache.OptSweepInterval(500 * time.Millisecond))
	go lc.Start()

	app.GET("/stats", func(r *web.Ctx) web.Result {
		return web.JSON.Result(lc.Stats())
	})

	app.GET("/item/:id", func(r *web.Ctx) web.Result {
		data, _, _ := lc.GetOrSet(
			web.StringValue(r.RouteParam("id")),
			getData,
			cache.OptValueTTL(30*time.Second),
			cache.OptValueOnRemove(func(key interface{}, reason cache.RemovalReason) {
				log.Infof("cache item removed: %s %v", key, reason)
			}),
		)
		return web.JSON.Result(data)
	})

	if err := graceful.Shutdown(app); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
