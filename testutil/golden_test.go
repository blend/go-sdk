/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/testutil"
	"github.com/blend/go-sdk/uuid"
)

const (
	// NOTE: This is "peace of mind" by obscurity; we assume callers will
	//       not use this flag, so we can register it with `MarkUpdateGoldenFlag()`
	//       and then reference it in `TestAssertGoldenFile_Assert()`.
	testUpdateGoldenFlag = "update-golden-0401b706-ae3c-4e7a-ba29-531fb811d3a2"
)

func TestGetTestFixture(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	data := testutil.GetTestFixture(it, "sentinel.txt")
	it.Equal("Any content will suffice here.\n", string(data))
}

func TestAssertGoldenFile_Assert(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	expected := strings.Join([]string{
		"#",
		"# Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved",
		"# Blend Confidential - Restricted",
		"#",
		"a:",
		"  b:",
		"    c: 10",
		"",
	}, "\n")
	testutil.AssertGoldenFile(it, []byte(expected), "config.yml", testutil.OptUpdateGoldenFlag(testUpdateGoldenFlag))
}

func TestAssertGoldenFile_Update(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	// Generate a test fixtures directory we control.
	tempDir, err := os.MkdirTemp("", "")
	it.Nil(err)
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	// Write a golden file to the test fixtures directory.
	path := filepath.Join(tempDir, "golden.focus.txt")
	err = os.WriteFile(path, []byte("short\n"), 0644)
	it.Nil(err)

	// NOTE: Add a per-run unique flag so that callers can't accidentally
	//       force an update.
	uniqueFlag := fmt.Sprintf("update-golden-%s", uuid.V4())

	// Run `AssertGoldenFile` first in `assert` mode so the `goldenFileFlags`
	// map can cache the flag (and then in subsequent runs we can enable
	// the flag).
	testutil.AssertGoldenFile(
		it, []byte("short\n"), "focus.txt",
		testutil.OptTestFixtureDirectory(tempDir),
		testutil.OptUpdateGoldenFlag(uniqueFlag),
	)

	// Now forcefully turn the flag on.
	f := flag.CommandLine.Lookup(uniqueFlag)
	it.NotNil(f)
	f.Value.Set("true")

	// Run in `update` mode (not `assert`)
	testutil.AssertGoldenFile(
		it, []byte("long-extra\n"), "focus.txt",
		testutil.OptTestFixtureDirectory(tempDir),
		testutil.OptUpdateGoldenFlag(uniqueFlag),
	)

	// Finally verify the file was overwritten.
	newContents, err := os.ReadFile(path)
	it.Nil(err)
	it.Equal([]byte("long-extra\n"), newContents)
}
