package configutil

import (
	"context"
	"io"

	"github.com/blend/go-sdk/env"
)

// ConfigOptions are options built for reading configs.
type ConfigOptions struct {
	Log         Logger
	Context     context.Context
	ContentsExt string
	Contents    io.Reader
	FilePaths   []string
	Env         env.Vars
}

// Background yields a context for a config options set.
func (co ConfigOptions) Background() context.Context {
	var background context.Context
	if co.Context != nil {
		background = co.Context
	} else {
		background = context.Background()
	}

	background = WithConfigFilePaths(background, co.FilePaths)
	background = WithEnvVars(background, co.Env)
	return background
}
