package configutil

import (
	"context"
	"io"

	"github.com/blend/go-sdk/env"
)

// Option is a modification of config options.
type Option func(*ConfigOptions) error

// OptLog sets the configutil logger.
func OptLog(log Logger) Option {
	return func(co *ConfigOptions) error {
		co.Log = log
		return nil
	}
}

// OptContext sets the context on the options.
func OptContext(ctx context.Context) Option {
	return func(co *ConfigOptions) error {
		co.Context = ctx
		return nil
	}
}

// OptContents sets the contents on the options.
func OptContents(ext string, contents io.Reader) Option {
	return func(co *ConfigOptions) error {
		co.ContentsExt = ext
		co.Contents = contents
		return nil
	}
}

// OptAddFilePaths adds paths to search for the config file.
func OptAddFilePaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = append(co.FilePaths, paths...)
		return nil
	}
}

// OptAddPreferredFilePaths adds paths to search first for the config file.
func OptAddPreferredFilePaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = append(paths, co.FilePaths...)
		return nil
	}
}

// OptFilePaths sets paths to search for the config file.
func OptFilePaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = paths
		return nil
	}
}

// OptEnv sets the config options environment variables.
// If unset, will default to the current global environment variables.
func OptEnv(vars env.Vars) Option {
	return func(co *ConfigOptions) error {
		co.Env = vars
		return nil
	}
}
