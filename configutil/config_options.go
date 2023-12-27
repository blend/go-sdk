/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"io"

	"github.com/blend/go-sdk/env"
)

// ConfigOptions are options built for reading configs.
type ConfigOptions struct {
	Log		Logger
	Context		context.Context
	Contents	[]ConfigContents
	FilePaths	[]string
	Env		env.Vars
}

// ConfigContents are literal contents to read from.
type ConfigContents struct {
	Ext		string
	Contents	io.Reader
}

// Background yields a context for a config options set.
func (co ConfigOptions) Background() context.Context {
	var background context.Context
	if co.Context != nil {
		background = co.Context
	} else {
		background = context.Background()
	}

	background = WithConfigPaths(background, co.FilePaths)
	background = env.WithVars(background, co.Env)
	return background
}
