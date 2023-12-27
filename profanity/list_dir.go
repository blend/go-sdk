/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"os"
)

// ListDir reads the directory named by dirname and returns
// a sorted list of directory entries.
func ListDir(path string) (dirs []os.FileInfo, files []os.FileInfo, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var children []os.FileInfo
	children, err = f.Readdir(-1)
	if err != nil {
		return
	}
	for _, child := range children {
		if child.IsDir() {
			dirs = append(dirs, child)
		} else {
			files = append(files, child)
		}
	}
	return
}
