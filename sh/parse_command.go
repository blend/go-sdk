/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sh

import (
	"github.com/blend/go-sdk/stringutil"
)

// ParseCommand returns the bin and args for a given statement.
//
// Typical use cases are to break apart single strings that represent
// comamnds and their arguments.
//
// Example:
//
//	bin, args := sh.ParseCommand("git rebase master")
//
// Would yield "git", and "rebase","master" as the args.
func ParseCommand(statement string) (bin string, args []string) {
	parts := stringutil.SplitSpaceQuoted(statement)
	if len(parts) > 0 {
		bin = parts[0]
	}
	if len(parts) > 1 {
		args = parts[1:]
	}
	return
}
