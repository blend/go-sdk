/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
)

// MustRead reads a config from optional path(s) and panics on error.
//
// It is functionally equivalent to `Read` outside error handling; see this function for more information.
func MustRead(ref Any, options ...Option) (filePaths []string) {
	var err error
	filePaths, err = Read(ref, options...)
	if !IsIgnored(err) {
		panic(err)
	}
	return
}

// Read reads a config from optional path(s), returning the paths read from (in the order visited), and an error if there were any issues.
/*
If the ref type is a `Resolver` the `Resolve(context.Context) error` method will
be called on the ref and passed a context configured from the given options.

By default, a well known set of paths will be read from (including a path read from the environment variable `CONFIG_PATH`).

You can override this by providing options to specify which paths will be read from:

	paths, err := configutil.Read(&cfg, configutil.OptPaths("foo.yml"))

The above will _only_ read from `foo.yml` to populate the `cfg` reference.
*/
func Read(ref Any, options ...Option) (paths []string, err error) {
	var configOptions ConfigOptions
	configOptions, err = createConfigOptions(options...)
	if err != nil {
		return
	}

	for _, contents := range configOptions.Contents {
		MaybeDebugf(configOptions.Log, "reading config contents with extension `%s`", contents.Ext)
		err = deserialize(contents.Ext, contents.Contents, ref)
		if err != nil {
			return
		}
	}

	var f *os.File
	var path string
	var resolveErr error
	for _, path = range configOptions.FilePaths {
		if path == "" {
			continue
		}
		MaybeDebugf(configOptions.Log, "checking for config path: %s", path)
		f, resolveErr = os.Open(path)
		if IsNotExist(resolveErr) {
			continue
		}
		if resolveErr != nil {
			err = ex.New(resolveErr)
			break
		}
		defer f.Close()

		MaybeDebugf(configOptions.Log, "reading config path: %s", path)
		resolveErr = deserialize(filepath.Ext(path), f, ref)
		if resolveErr != nil {
			err = ex.New(resolveErr)
			return
		}

		paths = append(paths, path)
	}

	if typed, ok := ref.(Resolver); ok {
		MaybeDebugf(configOptions.Log, "calling config resolver")
		if resolveErr := typed.Resolve(configOptions.Background()); resolveErr != nil {
			MaybeErrorf(configOptions.Log, "calling resolver error: %+v", resolveErr)
			err = resolveErr
			return
		}
	}
	return
}

func createConfigOptions(options ...Option) (configOptions ConfigOptions, err error) {
	configOptions.Env = env.Env()
	configOptions.FilePaths = DefaultPaths
	if configOptions.Env.Has(EnvVarConfigPath) {
		configOptions.FilePaths = append(configOptions.Env.CSV(EnvVarConfigPath), configOptions.FilePaths...)
	}
	for _, option := range options {
		if err = option(&configOptions); err != nil {
			return
		}
	}
	return
}

// deserialize deserializes a config.
func deserialize(ext string, r io.Reader, ref Any) error {
	// make sure the extension starts with a "."
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// based off the extension, use the appropriate deserializer
	switch strings.ToLower(ext) {
	case ExtensionJSON:
		return ex.New(json.NewDecoder(r).Decode(ref))
	case ExtensionYAML, ExtensionYML:
		return ex.New(yaml.NewDecoder(r).Decode(ref))
	default:	// return an error if we're passed a weird extension
		return ex.New(ErrInvalidConfigExtension, ex.OptMessagef("extension: %s", ext))
	}
}
