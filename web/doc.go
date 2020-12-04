/*
Package web implements a model view controller system for building http servers.
It is meant to be composed with other packages to form everything from small api servers
to fully formed web view applications.

Basics

To create a web server:

	import "golang.blend.com/sdk/graceful"
	import "golang.blend.com/sdk/logger"
	import "golang.blend.com/sdk/web"

	...

	app := web.MustNew(web.OptBindAddr(os.Getenv("BIND_ADDR")))
	app.GET("/", func(_ *web.Ctx) web.Result {
		return web.Text.Result("hello world")
	})
	if err := graceful.Shutdown(app); err != nil {
		logger.FatalExit(err)
	}

This will start a web server with a trivial endpoint mounted at the path "/" for the HTTP Verb "GET".
This example will also start the server and listen for SIGINT and SIGTERM os signals,
and close the server gracefully if they're recieved, letting requests finish.

There are many more examples in the golang.blend.com/sdk/examples/web directory.
*/
package web // import "github.com/blend/go-sdk/web"
