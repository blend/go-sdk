package configutil

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/yaml"
)

const (
	// EnvVarConfigPath is the env var for configs.
	EnvVarConfigPath = "CONFIG_PATH"
	// ExtensionJSON is a file extension.
	ExtensionJSON = ".json"
	// ExtensionYAML is a file extension.
	ExtensionYAML = ".yaml"
	// ExtensionYML is a file extension.
	ExtensionYML = ".yml"
)

var (
	// DefaultPaths are default path locations.
	// They are tested and read in order, so the later
	// paths will override data found in the earlier ones.
	DefaultPaths = []string{
		"/var/secrets/config.yml",
		"/var/secrets/config.yaml",
		"/var/secrets/config.json",
		"./_config/config.yml",
		"./_config/config.yaml",
		"./_config/config.json",
		"./config.yml",
		"./config.yaml",
		"./config.json",
	}
)

// Read reads a config from optional path(s).
// Paths will be tested from a standard set of defaults (ex. config.yml)
// and optionally a csv named in the `CONFIG_PATH` environment variable.
func Read(ref Any, options ...Option) (path string, err error) {
	var configOptions ConfigOptions
	configOptions, err = createConfigOptions(options...)
	if err != nil {
		return
	}

	if configOptions.Contents != nil {
		MaybeDebugf(configOptions.Log, "reading reader contents with extension `%s`", configOptions.ContentsExt)
		err = deserialize(configOptions.ContentsExt, configOptions.Contents, ref)
		if err != nil {
			return
		}
	} else {
		// for each of the paths
		// if the path doesn't exist, continue, read the path that is found.
		var f *os.File
		for _, path = range configOptions.FilePaths {
			if path == "" {
				continue
			}
			MaybeDebugf(configOptions.Log, "checking for config file %s", path)
			f, err = os.Open(path)
			if IsNotExist(err) {
				continue
			}
			if err != nil {
				err = ex.New(err)
				break
			}
			defer f.Close()
			MaybeDebugf(configOptions.Log, "reading config file %s", path)
			err = deserialize(filepath.Ext(path), f, ref)
			break
		}
		if err != nil && !IsNotExist(err) {
			return
		}
	}

	if typed, ok := ref.(Resolver); ok {
		MaybeDebugf(configOptions.Log, "calling config resolver")
		if resolveErr := typed.Resolve(configOptions.Background()); resolveErr != nil {
			err = resolveErr
			return
		}
	}
	return
}

func createConfigOptions(options ...Option) (configOptions ConfigOptions, err error) {
	configOptions.Env = env.Env()
	configOptions.FilePaths = DefaultPaths
	for _, option := range options {
		if err = option(&configOptions); err != nil {
			return
		}
	}
	if configOptions.Env.Has(EnvVarConfigPath) {
		configOptions.FilePaths = append(configOptions.Env.CSV(EnvVarConfigPath), configOptions.FilePaths...)
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
	default: // return an error if we're passed a weird extension
		return ex.New(ErrInvalidConfigExtension, ex.OptMessagef("extension: %s", ext))
	}
}
