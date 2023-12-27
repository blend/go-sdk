/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil

import (
	"os"

	"github.com/blend/go-sdk/ex"
)

// RemoveMany removes an array of files.
func RemoveMany(filePaths ...string) error {
	var err error
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			return ex.New(err, ex.OptMessage(path))
		}
	}
	return err
}
