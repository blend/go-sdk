/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// ShimTracer is a basic tracer.
type ShimTracer struct {
	OnStart		func(*web.Ctx)
	OnFinish	func(*web.Ctx, error)
}

// Start begins the trace.
func (st ShimTracer) Start(ctx *web.Ctx) web.TraceFinisher {
	if st.OnStart != nil {
		st.OnStart(ctx)
	}
	return &ShimTraceFinisher{parent: &st}
}

// ShimTraceFinisher finishes the traces.
type ShimTraceFinisher struct {
	parent *ShimTracer
}

// Finish closes the trace.
func (stf ShimTraceFinisher) Finish(ctx *web.Ctx, err error) {
	stf.parent.OnFinish(ctx, err)
}

func main() {
	log := logger.All()

	app := web.MustNew(
		web.OptBindAddr(":8080"),
		web.OptLog(log),
		web.OptTracer(ShimTracer{
			OnStart:	func(_ *web.Ctx) { log.Infof("Trace Started") },
			OnFinish:	func(_ *web.Ctx, _ error) { log.Infof("Trace Finished") },
		}),
	)
	app.GET("/", func(r *web.Ctx) web.Result {
		return web.Text.Result("ok!")
	})

	if err := graceful.Shutdown(app); err != nil {
		logger.FatalExit(err)
	}
}
