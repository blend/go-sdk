/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sourceutil

import (
	"context"
	"strings"
)

// StringSubstitution is a mutator for a string.
// It returns the modified string, and a bool if the rule matched or not.
type StringSubstitution func(context.Context, string) (string, bool)

// SubstituteString rewrites a string literal.
func SubstituteString(before, after string) StringSubstitution {
	return func(ctx context.Context, contents string) (string, bool) {
		if !strings.Contains(contents, before) {
			return "", false
		}
		return strings.Replace(contents, before, after, -1), true
	}
}
