/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"encoding/csv"
	"strings"
)

// CSV produces a csv from a given set of values.
// Deprecated: Use `encoding/csv.Writer` directly instead.
func CSV(values []string) string {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write(values); err != nil {
		return ""
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return ""
	}
	return strings.TrimSpace(builder.String())
}
