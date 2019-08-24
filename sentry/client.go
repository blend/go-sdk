package sentry

import (
	"net/http"
	"runtime"
	"time"

	raven "github.com/getsentry/sentry-go"

	ex "github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
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
			Dsn:         cfg.DSN,
			Environment: cfg.EnvironmentOrDefault(),
			ServerName:  cfg.ServerNameOrDefault(),
			Dist:        cfg.DistOrDefault(),
			Release:     cfg.ReleaseOrDefault(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		Config: cfg,
		Client: rc,
	}, nil
}

// Client is a wrapper for the sentry-go client.
type Client struct {
	Config Config
	Client *raven.Client
}

// Notify sends a notification.  This is primarily used for the logger.
func (c Client) Notify(ee *logger.ErrorEvent) {
	c.Client.CaptureEvent(errEvent(ee), nil, raven.NewScope())
	c.Client.Flush(time.Second)
}

// errEvent converts an `ErrorEvent` from the `Logger` package to a struct that
// can be consumed by the sentry SDK.
func errEvent(ee *logger.ErrorEvent) *raven.Event {
	return &raven.Event{
		Timestamp: ee.Timestamp().Unix(),
		Level:     raven.Level(ee.Flag()),
		Tags:      ee.Labels(),
		Extra:     errExtra(ee),
		Platform:  "go",
		Sdk: raven.SdkInfo{
			Name:    SDK,
			Version: raven.Version,
			Packages: []raven.SdkPackage{{
				Name:    SDK,
				Version: raven.Version,
			}},
		},
		Request: errRequest(ee),
		Message: ex.ErrClass(ee.Err()),
		Exception: []raven.Exception{
			{
				Type:       ex.ErrClass(ee.Err()),
				Value:      ex.ErrMessage(ee.Err()),
				Stacktrace: &raven.Stacktrace{Frames: errFrames(ee.Err())},
			},
		},
	}
}

// errExtra retrives annotations for an `ErrorEvent` and converts it to an
// interface map that can be consumed by raven
func errExtra(ee *logger.ErrorEvent) (annotations map[string]interface{}) {
	for k, v := range ee.Annotations() {
		annotations[k] = v
	}
	return
}

func errRequest(ee *logger.ErrorEvent) (requestMeta raven.Request) {
	if ee.State() == nil {
		return
	}
	typed, ok := ee.State().(*http.Request)
	if !ok {
		return
	}
	requestMeta = requestMeta.FromHTTPRequest(typed)
	return
}

func errFrames(err error) (output []raven.Frame) {
	stacktrace := ex.ErrStackTrace(err)
	if stacktrace == nil {
		return
	}
	pointers, ok := stacktrace.(ex.StackPointers)
	if !ok {
		return
	}

	runtimeFrames := runtime.CallersFrames(pointers)

	for {
		callerFrame, more := runtimeFrames.Next()
		output = append([]raven.Frame{
			raven.NewFrame(callerFrame),
		}, output...)
		if !more {
			break
		}
	}
	return
}
