/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package copyright

import (
	"os"
)

// Action is the action to run.
type Action func(path string, info os.FileInfo, file, notice []byte) error
