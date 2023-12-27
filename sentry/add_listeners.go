/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// AddListeners adds error listeners.
func AddListeners(log logger.Listenable, meta configmeta.Meta, cfg Config, opts ...AddListenersOption) error {
	if log == nil || cfg.IsZero() {
		return nil
	}
	if cfg.Environment == "" {
		cfg.Environment = meta.ServiceEnv
	}
	if cfg.ServerName == "" {
		cfg.ServerName = meta.ServiceName
	}
	if cfg.Release == "" {
		cfg.Release = meta.Version
	}

	if typed, ok := log.(logger.InfofReceiver); ok {
		typed.Infof("using sentry host: %s and environment: %s", webutil.MustParseURL(cfg.DSN).Hostname(), cfg.Environment)
	}
	client, err := New(cfg)
	if err != nil {
		return err
	}

	options := AddListenersOptions{
		EnabledFlags: DefaultListenerFlags,
		Scopes:       logger.ScopesAll(),
	}
	for _, opt := range opts {
		opt(&options)
	}

	listener := logger.NewScopedErrorEventListener(client.Notify, options.Scopes)
	for _, flag := range options.EnabledFlags {
		log.Listen(flag, ListenerName, listener)
	}
	return nil
}

// AddListenersOptions are all the options we can set when
// adding Sentry error listeners
type AddListenersOptions struct {
	EnabledFlags []string
	Scopes       *logger.Scopes
}

// AddListenersOption mutates AddListeners options
type AddListenersOption func(options *AddListenersOptions)

// AddListenersOptionFlags sets the logger flags to send Sentry
// notifications for
func AddListenersOptionFlags(flags ...string) AddListenersOption {
	return func(options *AddListenersOptions) {
		options.EnabledFlags = flags
	}
}

// AddListenersOptionScopes sets the logger scopes to send Sentry
// notifications for
func AddListenersOptionScopes(scopes ...string) AddListenersOption {
	return func(options *AddListenersOptions) {
		options.Scopes = logger.NewScopes(scopes...)
	}
}
