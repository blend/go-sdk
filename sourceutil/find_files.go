/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"context"
	"os"
	"path/filepath"
)

// FindFiles finds all files in a given path that matches a given glob
// but does not traverse recursively.
func FindFiles(ctx context.Context, targetPath string, matchGlob string) (output []string, err error) {
	err = filepath.Walk(targetPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if path == targetPath {
				return nil
			}
			return filepath.SkipDir
		}
		matched, err := filepath.Match(matchGlob, info.Name())
		if err != nil {
			return err
		}
		if matched {
			output = append(output, path)
		}
		return nil
	})
	return
}
