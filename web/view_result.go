/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"net/http"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// ViewResult is a result that renders a view.
type ViewResult struct {
	ViewName	string
	StatusCode	int
	ViewModel	interface{}
	Views		*ViewCache
	Template	*template.Template
}

// Render renders the result to the given response writer.
func (vr *ViewResult) Render(ctx *Ctx) (err error) {
	// you must set the template to be rendered.
	if vr.Template == nil {
		err = ex.New(ErrUnsetViewTemplate)
		return
	}

	if ctx.Tracer != nil {
		if typed, ok := ctx.Tracer.(ViewTracer); ok {
			tf := typed.StartView(ctx, vr)
			defer func() {
				tf.FinishView(ctx, vr, err)
			}()
		}
	}

	ctx.Response.Header().Set(webutil.HeaderContentType, webutil.ContentTypeHTML)

	// use a pooled buffer if possible
	var buffer *bytes.Buffer
	if vr.Views != nil && vr.Views.BufferPool != nil {
		buffer = vr.Views.BufferPool.Get()
		defer vr.Views.BufferPool.Put(buffer)
	} else {
		buffer = new(bytes.Buffer)
	}

	err = vr.Template.Execute(buffer, &ViewModel{
		Env:	env.Env(),
		Ctx:	ctx,
		Status: ViewStatus{
			Text:	http.StatusText(vr.StatusCode),
			Code:	vr.StatusCode,
		},
		ViewModel:	vr.ViewModel,
	})
	if err != nil {
		err = ex.New(err)
		if ctx.App != nil {
			logger.MaybeErrorContext(ctx.Context(), ctx.App.Log, err)
		}
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		_, _ = ctx.Response.Write([]byte(fmt.Sprintf("%+v\n", err)))
		return
	}

	ctx.Response.WriteHeader(vr.StatusCode)
	_, err = ctx.Response.Write(buffer.Bytes())
	if err != nil {
		if typed, ok := err.(*net.OpError); ok {
			err = ex.New(webutil.ErrNetWrite, ex.OptInner(typed))
			return
		}
		err = ex.New(err)
		return
	}
	return
}
