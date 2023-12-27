/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
