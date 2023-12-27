/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import "strings"

// TSV produces a tab seprated values from a given set of values.
func TSV(values []string) string {
	return strings.Join(values, "\t")
}
