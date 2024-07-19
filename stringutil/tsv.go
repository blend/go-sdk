/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import "strings"

// TSV produces a tab seprated values from a given set of values.
func TSV(values []string) string {
	return strings.Join(values, "\t")
}
