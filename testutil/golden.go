/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/blend/go-sdk/assert"
)

const (
	// DefaultFixtureDirectory is the default directory where test fixture
	// files will be retrieved. Since `go test` sets the present working directory
	// to the current directory under test, this is a relative path.
	DefaultFixtureDirectory = "./testdata"
	// DefaultGoldenFilePrefix is the prefix that will be prepended to suffixes
	// for golden filenames.
	DefaultGoldenFilePrefix = "golden."
	// DefaultUpdateGoldenFlag is the flag that this package will use to check
	// if golden files should be updated.
	DefaultUpdateGoldenFlag = "update-golden"

	defaultNewFilePermissions = 0644
)

var (
	goldenFileFlags     = map[string]*bool{}
	goldenFileFlagsLock = sync.Mutex{}
)

// GetTestFixture opens a file in the test fixtures directory. This
// relies on the present working directory being set to the current
// directory under test and will default to `./testdata` when
// reading files.
func GetTestFixture(it *assert.Assertions, filename string, opts ...FixtureOption) []byte {
	fc := NewFixtureConfig(opts...)
	path := filepath.Join(fc.Directory, filename)
	data, err := os.ReadFile(path)
	it.Nil(err, fmt.Sprintf("Failed reading Test Fixture File %q", path))
	return data
}

// AssertGoldenFile checks that a "golden file" matches the expected contents.
// The golden file will use the test fixtures directory and will use the
// filename suffix **after** a prefix of `golden.`, e.g.
// `./testdata/golden.config.yml`.
//
// Managing golden files can be tedious, so test suites can optionally specify
// a boolean flag (e.g. `go test --update-golden`) that can be propagated to
// this function; in which case the golden file will just be overwritten instead
// of compared against `expected`.
func AssertGoldenFile(it *assert.Assertions, expected []byte, filenameSuffix string, opts ...FixtureOption) {
	fc := NewFixtureConfig(opts...)
	update := getUpdateGoldenFlag(it, fc.UpdateGoldenFlag)
	filename := fmt.Sprintf("%s%s", fc.GoldenFilePrefix, filenameSuffix)
	path := filepath.Join(fc.Directory, filename)

	if update {
		// NOTE: `defaultNewFilePermissions` will only be used for **new**
		//       files, otherwise `os.WriteFile()` will preserve permissions.
		err := os.WriteFile(path, expected, defaultNewFilePermissions)
		it.Nil(err, fmt.Sprintf("Error writing Golden File %q", path))
		return
	}

	actual, err := os.ReadFile(path)
	it.Nil(err, fmt.Sprintf("Failed reading Golden File %q", path))
	it.True(bytes.Equal(expected, actual), fmt.Sprintf("Golden File %q does not match expected, consider running with the '--%s' flag to update the Golden File", path, fc.UpdateGoldenFlag))
}

// MarkUpdateGoldenFlag is intended to be used in a `TestMain()` to declare
// a flag **before** `go test` parses flags (if not, unknown flags will
// fail a test).
//
// This is expected to be used in two modes:
// > MarkUpdateGoldenFlag()
// which will mark `--update-golden` (via `DefaultUpdateGoldenFlag`) as a valid
// flag for tests and with a custom flag override:
// > MarkUpdateGoldenFlag(OptUpdateGoldenFlag(customFlag))
func MarkUpdateGoldenFlag(opts ...FixtureOption) {
	goldenFileFlagsLock.Lock()
	defer goldenFileFlagsLock.Unlock()

	fc := NewFixtureConfig(opts...)
	b := flag.Bool(fc.UpdateGoldenFlag, false, "Update Golden Files")
	goldenFileFlags[fc.UpdateGoldenFlag] = b
}

// getUpdateGoldenFlag returns a (parsed) command line flag that indicates if
// golden files should be updated.
//
// NOTE: This function is careful not to re-define a flag that has already
//
//	been set. This can also be checked in `flag.CommandLine.`
//	See: https://github.com/golang/go/blob/go1.17.1/src/flag/flag.go#L879
func getUpdateGoldenFlag(it *assert.Assertions, name string) bool {
	goldenFileFlagsLock.Lock()
	defer goldenFileFlagsLock.Unlock()

	b, ok := goldenFileFlags[name]
	// Careful not to re-define a flag that has already been set.
	if !ok {
		existing := flag.CommandLine.Lookup(name)
		it.Nil(existing, fmt.Sprintf("Update Golden Flag '--%s' is already defined elsewhere", name))
		b = flag.Bool(name, false, "Update Golden Files")
		goldenFileFlags[name] = b
	}
	flag.Parse()
	it.NotNil(b, fmt.Sprintf("Parsed boolean flag '--%s' should not be nil", name))
	return *b
}

// FixtureConfig represents defaults used for working with test fixtures.
type FixtureConfig struct {
	Directory        string
	GoldenFilePrefix string
	UpdateGoldenFlag string
}

// NewFixtureConfig returns a new `FixtureConfig` and applies options.
func NewFixtureConfig(opts ...FixtureOption) FixtureConfig {
	fc := FixtureConfig{
		Directory:        DefaultFixtureDirectory,
		GoldenFilePrefix: DefaultGoldenFilePrefix,
		UpdateGoldenFlag: DefaultUpdateGoldenFlag,
	}
	for _, opt := range opts {
		opt(&fc)
	}
	return fc
}

// FixtureOption is a mutator for a `FixtureConfig`.
type FixtureOption func(*FixtureConfig)

// OptTestFixtureDirectory sets the directory used to look up test fixture
// files. Both relative and absolute paths are supported.
func OptTestFixtureDirectory(name string) FixtureOption {
	return func(fc *FixtureConfig) {
		fc.Directory = name
	}
}

// OptUpdateGoldenFlag sets the default flag used to determine if golden files
// should be updated (e.g. to use `--update-golden`, pass `"update-golden"`
// here).
func OptUpdateGoldenFlag(name string) FixtureOption {
	return func(fc *FixtureConfig) {
		fc.UpdateGoldenFlag = name
	}
}
