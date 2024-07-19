/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"os"
	"path/filepath"
)

// FileInfo extends os.FileInfo with the full path.
type FileInfo struct {
	os.FileInfo
	FullPath string
}

// LS returns a list of files for a given path.
func LS(root string) (output []FileInfo, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		}
		output = append(output, FileInfo{
			FileInfo: info,
			FullPath: path,
		})
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	return
}
