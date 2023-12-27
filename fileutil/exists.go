/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil

import (
	"os"

	"github.com/blend/go-sdk/ex"
)

// pathInfo encapsulates two of the more salient outputs of `os.Stat()`
// - Can we stat the file / path? (i.e. does it exist?)
// - Is the file mode a directory?
type pathInfo struct {
	Exists	bool
	IsDir	bool
}

func pathExists(path string) (*pathInfo, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return &pathInfo{Exists: true, IsDir: fi.Mode().IsDir()}, nil
	}

	if os.IsNotExist(err) {
		return &pathInfo{Exists: false}, nil
	}

	return nil, ex.New(err)
}

// FileExists determines if a path exists and is a file (not a directory).
func FileExists(path string) (bool, error) {
	pi, err := pathExists(path)
	if err != nil {
		return false, err
	}

	if pi.Exists && pi.IsDir {
		return false, ex.New("Path exists but is a directory", ex.OptMessagef("Path: %q", path))
	}

	return pi.Exists, nil
}

// DirExists determines if a path exists and is a directory (not a file).
func DirExists(path string) (bool, error) {
	pi, err := pathExists(path)
	if err != nil {
		return false, err
	}

	if pi.Exists && !pi.IsDir {
		return false, ex.New("Path exists but is not a directory", ex.OptMessagef("Path: %q", path))
	}

	return pi.Exists, nil
}
