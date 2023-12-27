/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/fileutil"
)

func TestFileExists(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	// Exists: Yes, File: Yes
	f := tempFile(it)
	exists, err := fileutil.FileExists(f.Name())
	it.Nil(err)
	it.True(exists)

	// Exists: No
	name := f.Name() + "-missing"
	exists, err = fileutil.FileExists(name)
	it.Nil(err)
	it.False(exists)

	// Exists: Yes, File: No
	dir := tempDir(it)
	exists, err = fileutil.FileExists(dir)
	it.False(exists)
	it.NotNil(err)
	expected := fmt.Sprintf("Path exists but is a directory; Path: %q", dir)
	it.Equal(expected, fmt.Sprintf("%v", err))

	// Stat fails
	path := tooLongPath()
	exists, err = fileutil.FileExists(path)
	it.False(exists)
	it.NotNil(err)
	expected = fmt.Sprintf("stat %s: file name too long\nfile name too long", path)
	it.Equal(expected, fmt.Sprintf("%v", err))
}

func TestDirExists(t *testing.T) {
	t.Parallel()
	it := assert.New(t)

	// Exists: Yes, File: Yes
	f := tempFile(it)
	exists, err := fileutil.DirExists(f.Name())
	it.False(exists)
	it.NotNil(err)
	expected := fmt.Sprintf("Path exists but is not a directory; Path: %q", f.Name())
	it.Equal(expected, fmt.Sprintf("%v", err))

	// Exists: No
	name := f.Name() + "-missing"
	exists, err = fileutil.DirExists(name)
	it.Nil(err)
	it.False(exists)

	// Exists: Yes, File: No
	dir := tempDir(it)
	exists, err = fileutil.DirExists(dir)
	it.Nil(err)
	it.True(exists)

	// Stat fails
	path := tooLongPath()
	exists, err = fileutil.DirExists(path)
	it.False(exists)
	it.NotNil(err)
	expected = fmt.Sprintf("stat %s: file name too long\nfile name too long", path)
	it.Equal(expected, fmt.Sprintf("%v", err))
}

func tempFile(it *assert.Assertions) *os.File {
	f, err := os.CreateTemp("", "")
	it.Nil(err)

	it.T.Cleanup(func() {
		err := f.Close()
		it.Nil(err)
		err = os.Remove(f.Name())
		it.Nil(err)
	})

	return f
}

func tempDir(it *assert.Assertions) string {
	dir, err := os.MkdirTemp("", "")
	it.Nil(err)

	it.T.Cleanup(func() {
		err = os.RemoveAll(dir)
		it.Nil(err)
	})

	return dir
}

func tooLongPath() string {
	parts := make([]string, 1024)
	parts[0] = ""
	parts[1] = "path"
	parts[2] = "to"
	for i := 3; i < 1024; i++ {
		parts[i] = "more-content"
	}
	return strings.Join(parts, string(os.PathSeparator))
}
