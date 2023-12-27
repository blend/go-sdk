/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"context"
	"time"
)

type appKey struct{}

// WithApp adds an app to a context.
func WithApp(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appKey{}, app)
}

// GetApp gets an app off a context.
func GetApp(ctx context.Context) *App {
	if value := ctx.Value(appKey{}); value != nil {
		if typed, ok := value.(*App); ok {
			return typed
		}
	}
	return nil
}

type requestStartedKey struct{}

// WithRequestStarted sets the request started time on a context.
func WithRequestStarted(ctx context.Context, requestStarted time.Time) context.Context {
	return context.WithValue(ctx, requestStartedKey{}, requestStarted)
}

// GetRequestStarted gets the request started time from a context.
func GetRequestStarted(ctx context.Context) time.Time {
	if value := ctx.Value(requestStartedKey{}); value != nil {
		if typed, ok := value.(time.Time); ok {
			return typed
		}
	}
	return time.Time{}
}

type sessionKey struct{}

// WithSession sets a session on a context.
func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, session)
}

// GetSession gets a session off a context.
func GetSession(ctx context.Context) *Session {
	if value := ctx.Value(sessionKey{}); value != nil {
		if typed, ok := value.(*Session); ok {
			return typed
		}
	}
	return nil
}
