/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"log"
	"sync"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// Any is a type alias for interface{}
type Any = interface{}

// APIController implements a simple controller.
type APIController struct {
	sync.Mutex
	db map[string]Any
}

// Register adds routes for the controller to the app.
func (ac *APIController) Register(app *web.App) {
	app.GET("/", ac.index)
	app.GET("/api", ac.all)
	app.GET("/api/:key", ac.get)
	app.POST("/api/:key", ac.post)
	app.PUT("/api/:key", ac.put)
	app.DELETE("/api/:key", ac.delete)
}

func (ac *APIController) ensureDB() {
	if ac.db == nil {
		ac.db = map[string]Any{}
	}
}

func (ac *APIController) index(r *web.Ctx) web.Result {
	return web.JSON.OK()
}

func (ac *APIController) all(r *web.Ctx) web.Result {
	ac.Lock()
	defer ac.Unlock()
	ac.ensureDB()

	return web.JSON.Result(ac.db)
}

func (ac *APIController) get(r *web.Ctx) web.Result {
	ac.Lock()
	defer ac.Unlock()
	ac.ensureDB()

	key, err := r.Param("key")
	if err != nil {
		return web.JSON.BadRequest(err)
	}

	value, hasValue := ac.db[key]
	if !hasValue {
		return web.JSON.NotFound()
	}
	return web.JSON.Result(value)
}

func (ac *APIController) post(r *web.Ctx) web.Result {
	ac.Lock()
	defer ac.Unlock()
	ac.ensureDB()

	body, err := r.PostBody()
	if err != nil {
		return web.JSON.InternalError(err)
	}

	key, err := r.Param("key")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	ac.db[key] = string(body)
	return web.JSON.OK()
}

func (ac *APIController) put(r *web.Ctx) web.Result {
	ac.Lock()
	defer ac.Unlock()
	ac.ensureDB()

	key, err := r.Param("key")
	if err != nil {
		return web.JSON.BadRequest(err)
	}

	_, hasValue := ac.db[key]
	if !hasValue {
		return web.JSON.NotFound()
	}

	body, err := r.PostBody()
	if err != nil {
		return web.JSON.InternalError(err)
	}
	ac.db[key] = string(body)

	return web.JSON.OK()
}

func (ac *APIController) delete(r *web.Ctx) web.Result {
	ac.Lock()
	defer ac.Unlock()
	ac.ensureDB()

	key, err := r.Param("key")
	if err != nil {
		return web.JSON.BadRequest(err)
	}

	_, hasValue := ac.db[key]
	if !hasValue {
		return web.JSON.NotFound()
	}
	delete(ac.db, key)
	return web.JSON.OK()
}

func main() {
	app := web.MustNew(web.OptLog(logger.All()))
	app.Register(new(APIController))
	if err := graceful.Shutdown(app); err != nil {
		log.Fatal(err)
	}
}
