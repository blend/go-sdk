/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"bytes"
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

// OptContents sets config contents on the options.
func OptContents(contents ...ConfigContents) Option {
	return func(co *ConfigOptions) error {
		co.Contents = contents
		return nil
	}
}

// OptAddContent adds contents to the options as a reader.
func OptAddContent(ext string, content io.Reader) Option {
	return func(co *ConfigOptions) error {
		co.Contents = append(co.Contents, ConfigContents{
			Ext:      ext,
			Contents: content,
		})
		return nil
	}
}

// OptAddContentString adds contents to the options as a string.
func OptAddContentString(ext string, contents string) Option {
	return func(co *ConfigOptions) error {
		co.Contents = append(co.Contents, ConfigContents{
			Ext:      ext,
			Contents: bytes.NewReader([]byte(contents)),
		})
		return nil
	}
}

// OptAddPaths adds paths to search for the config file.
//
// These paths will be added after the default paths.
func OptAddPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = append(co.FilePaths, paths...)
		return nil
	}
}

// OptAddFilePaths is deprecated; use `OptAddPaths`
func OptAddFilePaths(paths ...string) Option {
	return OptAddPaths(paths...)
}

// OptAddPreferredPaths adds paths to search first for the config file.
func OptAddPreferredPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = append(paths, co.FilePaths...)
		return nil
	}
}

// OptPaths sets paths to search for the config file.
func OptPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = paths
		return nil
	}
}

// OptUnsetPaths removes default paths from the paths set.
func OptUnsetPaths() Option {
	return func(co *ConfigOptions) error {
		co.FilePaths = nil
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
