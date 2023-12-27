/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package codeowners

import (
	"path/filepath"
	"strings"
)

// MakeRepositoryAbsolute make a path absolute.
func MakeRepositoryAbsolute(repositoryRoot, path string) (string, error) {
	var err error
	if !filepath.IsAbs(repositoryRoot) {
		repositoryRoot, err = filepath.Abs(repositoryRoot)
		if err != nil {
			return "", err
		}
	}
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}
	path = strings.TrimPrefix(path, repositoryRoot)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path, nil
}
