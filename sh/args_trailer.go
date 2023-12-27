/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrFlagsNoTrailer ex.Class = "sh; error parsing args trailer; missing '--' token, or nothing follows it"
)

// ArgsTrailer parses a set of os.Args, and returns everything after the `--` token.
// If there is no `--` token, an exception class "ErrFlagsNoTrailer" is returned.
func ArgsTrailer(args ...string) ([]string, error) {
	var foundIndex int
	for index, arg := range args {
		if strings.TrimSpace(arg) == "--" {
			foundIndex = index
			break
		}
	}
	if foundIndex == 0 {
		return nil, ex.New(ErrFlagsNoTrailer, ex.OptMessagef("args: %v", strings.Join(args, " ")))
	}
	if foundIndex == len(args)-1 {
		return nil, ex.New(ErrFlagsNoTrailer, ex.OptMessage("cannot be the last flag argument"))
	}
	return args[foundIndex+1:], nil
}
