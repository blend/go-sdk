/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

// SplitColon splits a string on a **single** colon. For example, for a basic
// auth header, we'd want to split a string of the form "<username>:<password>".
func SplitColon(value string) (first, second string, err error) {
	pair := strings.SplitN(value, ":", 2)
	if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
		err = ex.New(ErrInvalidSplitColonInput, ex.OptMessagef("input: %q", value))
		return
	}
	first = pair[0]
	second = pair[1]
	return
}
