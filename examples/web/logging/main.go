/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"
)

func main() {
	log := logger.Prod()
	app := web.MustNew(web.OptLog(log))
	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Text.Result("foo")
	})
	log.Listen(webutil.FlagHTTPRequest, logger.DefaultListenerName, webutil.NewHTTPRequestEventListener(func(_ context.Context, wre webutil.HTTPRequestEvent) {
		log.Infof("got a new request at route: %s", wre.Route)
	}))

	graceful.Shutdown(app)
}
