/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

var (
	_ Client = (*HTTPClient)(nil)
)

// HTTPClient is an implementation of the http client.
type HTTPClient struct {
	Config   Config
	Defaults []r2.Option
}

// Request creates a request with a context and a given set of options.
func (hc HTTPClient) Request(ctx context.Context, opts ...r2.Option) *r2.Request {
	callOptions := []r2.Option{
		r2.OptContext(ctx),
		r2.OptHeaderValue(webutil.HeaderContentType, webutil.ContentTypeApplicationJSON),
		r2.OptHeaderValue(webutil.HeaderAccept, "application/vnd.pagerduty+json;version=2"),
		r2.OptHeaderValue(webutil.HeaderAuthorization, fmt.Sprintf("Token token=%s", hc.Config.Token)),
	}
	if hc.Config.Email != "" {
		callOptions = append(callOptions,
			r2.OptHeaderValue(http.CanonicalHeaderKey("From"), hc.Config.Email),
		)
	}
	baseOptions := append(hc.Defaults,
		callOptions...,
	)
	return r2.New(hc.Config.Addr, append(baseOptions, opts...)...)
}
