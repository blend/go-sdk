/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"os"
)

// Action is the action to run.
type Action func(path string, info os.FileInfo, file, notice []byte) error
