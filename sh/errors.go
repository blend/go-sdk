/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"os/exec"
	"syscall"
)

// IsEPIPE is the epipe erorr.
func IsEPIPE(err error) bool {
	if typed, ok := err.(*exec.ExitError); ok {
		status := typed.Sys().(syscall.WaitStatus)
		if status.Signal() == syscall.SIGPIPE {
			return true
		}
	}
	return false
}
