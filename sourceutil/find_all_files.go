/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// FindAllFiles finds all the files that match a given glob recursively.
func FindAllFiles(ctx context.Context, startPath, matchGlob string) (output []string, err error) {
	err = filepath.Walk(startPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if path == startPath {
				return nil
			}
			if strings.HasPrefix(info.Name(), "_") {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			if strings.HasPrefix(path, "vendor/") {
				return filepath.SkipDir
			}
			return nil
		}

		matched, err := filepath.Match(matchGlob, info.Name())
		if err != nil {
			return err
		}
		if matched {
			if !strings.HasPrefix(path, "./") {
				path = "./" + path
			}
			output = append(output, path)
		}
		return nil
	})
	return
}
