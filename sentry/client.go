/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sentry

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	raven "github.com/blend/sentry-go"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

var (
	_ Sender = (*Client)(nil)
)

// MustNew returns a new client and panics on error.
func MustNew(cfg Config) *Client {
	c, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return c
}

// New returns a new client.
func New(cfg Config) (*Client, error) {
	rc, err := raven.NewClient(
		raven.ClientOptions{
			Dsn:		cfg.DSN,
			Environment:	cfg.Environment,
			ServerName:	cfg.ServerName,
			Dist:		cfg.Dist,
			Release:	cfg.Release,
			Debug:		cfg.Debug,
			DebugWriter:	os.Stderr,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		Config:	cfg,
		Client:	rc,
	}, nil
}

// Client is a wrapper for the sentry-go client.
type Client struct {
	Config	Config
	Client	*raven.Client
}

// Notify sends a notification.
func (c Client) Notify(ctx context.Context, ee logger.ErrorEvent) {
	c.Client.CaptureEvent(errEvent(ctx, ee), nil, raven.NewScope())
	c.Client.Flush(c.Config.FlushTimeoutOrDefault())	// goose this a bit
}

func errEvent(ctx context.Context, ee logger.ErrorEvent) *raven.Event {
	exceptions := []raven.Exception{
		{
			Type:		ex.ErrClass(ee.Err).Error(),
			Value:		ex.ErrMessage(ee.Err),
			Stacktrace:	errStackTrace(ee.Err),
		},
	}

	var innerErr error
	for innerErr = ex.ErrInner(ee.Err); innerErr != nil; innerErr = ex.ErrInner(innerErr) {
		rex := raven.Exception{
			Type:	ex.ErrClass(innerErr).Error(),
			Value:	ex.ErrMessage(innerErr),
		}
		if st := errStackTrace(innerErr); st != nil && len(st.Frames) > 0 {
			rex.Stacktrace = st
		}
		exceptions = append(exceptions, rex)
	}

	event := &raven.Event{
		Title:		ex.ErrClass(ee.Err).Error(),
		Timestamp:	logger.GetEventTimestamp(ctx, ee),
		Fingerprint:	errFingerprint(ctx, ex.ErrClass(ee.Err).Error()),
		Level:		raven.Level(ee.GetFlag()),
		Tags:		errTags(ctx),
		Extra:		errExtra(ctx),
		Platform:	"go",
		Contexts: map[string]interface{}{
			"device": map[string]interface{}{
				"arch":		runtime.GOARCH,
				"num_cpu":	runtime.NumCPU(),
			},
			"os": map[string]interface{}{
				"name": runtime.GOOS,
			},
			"runtime": map[string]interface{}{
				"name":			"go",
				"version":		runtime.Version(),
				"go_numroutines":	runtime.NumGoroutine(),
				"go_maxprocs":		runtime.GOMAXPROCS(0),
				"go_numcgocalls":	runtime.NumCgoCall(),
			},
		},
		Sdk: raven.SdkInfo{
			Name:		SDK,
			Version:	raven.Version,
			Packages: []raven.SdkPackage{{
				Name:		SDK,
				Version:	raven.Version,
			}},
		},
		Request:	errRequest(ee),
		Message:	ex.ErrClass(ee.Err).Error(),
		Exception:	exceptions,
	}

	// Set contextual information preserving existing data. For each context, if
	// the existing value is not of type map[string]interface{}, then no
	// additional information is added.
	if deviceContext, ok := event.Contexts["device"].(map[string]interface{}); ok {
		if _, ok := deviceContext["arch"]; !ok {
			deviceContext["arch"] = runtime.GOARCH
		}
		if _, ok := deviceContext["num_cpu"]; !ok {
			deviceContext["num_cpu"] = runtime.NumCPU()
		}
	}
	if osContext, ok := event.Contexts["os"].(map[string]interface{}); ok {
		if _, ok := osContext["name"]; !ok {
			osContext["name"] = runtime.GOOS
		}
	}
	if runtimeContext, ok := event.Contexts["runtime"].(map[string]interface{}); ok {
		if _, ok := runtimeContext["name"]; !ok {
			runtimeContext["name"] = "go"
		}
		if _, ok := runtimeContext["version"]; !ok {
			runtimeContext["version"] = runtime.Version()
		}
		if _, ok := runtimeContext["go_numroutines"]; !ok {
			runtimeContext["go_numroutines"] = runtime.NumGoroutine()
		}
		if _, ok := runtimeContext["go_maxprocs"]; !ok {
			runtimeContext["go_maxprocs"] = runtime.GOMAXPROCS(0)
		}
		if _, ok := runtimeContext["go_numcgocalls"]; !ok {
			runtimeContext["go_numcgocalls"] = runtime.NumCgoCall()
		}
	}

	return event
}

func errFingerprint(ctx context.Context, extra ...string) []string {
	if fingerprint := GetFingerprint(ctx); fingerprint != nil {
		return fingerprint
	}
	return extra
}

func errTags(ctx context.Context) map[string]string {
	labels := logger.GetLabels(ctx)
	if labels == nil {
		labels = make(map[string]string)
	}
	if hostname := env.GetVars(ctx).Hostname(); hostname != "" {
		labels["hostname"] = hostname
	}
	return labels
}

func errExtra(ctx context.Context) map[string]interface{} {
	return logger.GetAnnotations(ctx)
}

func errRequest(ee logger.ErrorEvent) *raven.Request {
	if ee.State == nil {
		return new(raven.Request)
	}
	typed, ok := ee.State.(*http.Request)
	if !ok {
		return &raven.Request{}
	}

	return newRavenRequest(typed)
}

func newRavenRequest(r *http.Request) *raven.Request {
	protocol := webutil.SchemeHTTP
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = webutil.SchemeHTTPS
	}
	url := fmt.Sprintf("%s://%s%s", protocol, r.Host, r.URL.Path)
	headers := make(map[string]string)
	headers["Host"] = r.Host
	var env map[string]string
	if addr, port, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		env = map[string]string{"REMOTE_ADDR": addr, "REMOTE_PORT": port}
	}
	return &raven.Request{
		URL:		url,
		Method:		r.Method,
		QueryString:	r.URL.RawQuery,
		Headers:	headers,
		Env:		env,
	}
}

func errStackTrace(err error) *raven.Stacktrace {
	if err != nil {
		frames := errFrames(err)
		if len(frames) > 0 {
			return &raven.Stacktrace{
				Frames: errFrames(err),
			}
		}
		return nil
	}
	return nil
}

func errFrames(err error) []raven.Frame {
	stacktrace := ex.ErrStackTrace(err)
	if stacktrace == nil {
		return nil
	}
	pointers, ok := stacktrace.(ex.StackPointers)
	if !ok {
		return nil
	}

	var output []raven.Frame
	runtimeFrames := runtime.CallersFrames(pointers)
	for {
		callerFrame, more := runtimeFrames.Next()
		// append in reverse order ...
		output = append([]raven.Frame{
			raven.NewFrame(callerFrame),
		}, output...)
		if !more {
			break
		}
	}
	return output
}
