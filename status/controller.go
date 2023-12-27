/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// NewController returns a new controller
func NewController(opts ...ControllerOption) *Controller {
	controller := Controller{
		Checks:		NewFreeform(),
		TrackedActions:	NewTrackedActionAggregator(),
	}
	for _, opt := range opts {
		opt(&controller)
	}
	return &controller
}

// ControllerOption mutates a controller.
type ControllerOption func(c *Controller)

// OptConfig returns an option that sets the configmeta.
func OptConfig(cfg configmeta.Meta) ControllerOption {
	return func(c *Controller) {
		c.Config = cfg
	}
}

// OptLog returns an option that sets the logger.
func OptLog(log logger.Log) ControllerOption {
	return func(c *Controller) {
		c.Checks.Log = log
		c.TrackedActions.Log = log
	}
}

// OptTimeout returns an option that sets checks timeout.
func OptTimeout(timeout time.Duration) ControllerOption {
	return func(c *Controller) {
		c.Checks.Timeout = timeout
	}
}

// OptUseBareMethods returns an option that sets if we should use bare methods.
func OptUseBareMethods(useBareMethods bool) ControllerOption {
	return func(c *Controller) {
		c.UseBareMethods = useBareMethods
	}
}

// OptCheck adds an sla check for a given service.
func OptCheck(serviceName string, check Checker) ControllerOption {
	return func(c *Controller) {
		if c.Checks.ServiceChecks == nil {
			c.Checks.ServiceChecks = make(map[string]Checker)
		}
		c.Checks.ServiceChecks[serviceName] = check
	}
}

// OptMiddleware adds default middleware for the status routes.
//
// Middleware must be set _before_ you register the controller.
func OptMiddleware(middleware ...web.Middleware) ControllerOption {
	return func(c *Controller) {
		c.Middleware = append(c.Middleware, middleware...)
	}
}

// Controller is a handler for the status endpoint.
//
// It will register `/status/sla` and `/status/details` routes
// on the given app.
type Controller struct {
	Config		configmeta.Meta
	Checks		*Freeform
	TrackedActions	*TrackedActionAggregator
	Middleware	[]web.Middleware
	UseBareMethods	bool
}

// Register adds the controller's routes to the app.
func (c Controller) Register(app *web.App) {
	if c.UseBareMethods {
		app.MethodBare(http.MethodGet, "/status", c.getStatus, c.Middleware...)
		app.MethodBare(http.MethodGet, "/status/sla", c.getStatusSLA, c.Middleware...)
		app.MethodBare(http.MethodGet, "/status/details", c.getStatusDetails, c.Middleware...)
	} else {
		app.GET("/status", c.getStatus, c.Middleware...)
		app.GET("/status/sla", c.getStatusSLA, c.Middleware...)
		app.GET("/status/details", c.getStatusDetails, c.Middleware...)
	}
}

// Interceptor returns a new interceptor for a given serviceName.
func (c Controller) Interceptor(serviceName string) async.Interceptor {
	return c.TrackedActions.Interceptor(serviceName)
}

// GET /status
func (c Controller) getStatus(r *web.Ctx) web.Result {
	return web.JSON.Result(c.Config)
}

// GET /status/sla
func (c Controller) getStatusSLA(r *web.Ctx) web.Result {
	return c.Checks.Endpoint()(r)
}

// GET /status/details
func (c Controller) getStatusDetails(r *web.Ctx) web.Result {
	return c.TrackedActions.Endpoint()(r)
}
