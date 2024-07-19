/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"encoding/csv"
	"strings"
)

// SplitCSV splits a corpus by the `,`.
// Deprecated: Use `encoding/csv.Reader` directly instead.
func SplitCSV(text string) []string {
	if len(text) == 0 {
		return nil
	}
	reader := csv.NewReader(strings.NewReader(text))
	output, err := reader.Read()
	if err != nil {
		return nil
	}
	return output
}
