/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package diff

import "bytes"

// Text converts a []Diff into a text report.
func Text(diffs []Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case DiffInsert:
			_, _ = buff.WriteString("+")
			_, _ = buff.WriteString(text)
		case DiffDelete:
			_, _ = buff.WriteString("-")
			_, _ = buff.WriteString(text)
		}
	}
	return buff.String()
}

// Text1 computes and returns the source text (all equalities and deletions).
func Text1(diffs []Diff) string {
	//StringBuilder text = new StringBuilder()
	var text bytes.Buffer

	for _, aDiff := range diffs {
		if aDiff.Type != DiffInsert {
			_, _ = text.WriteString(aDiff.Text)
		}
	}
	return text.String()
}

// Text2 computes and returns the destination text (all equalities and insertions).
func Text2(diffs []Diff) string {
	var text bytes.Buffer

	for _, aDiff := range diffs {
		if aDiff.Type != DiffDelete {
			_, _ = text.WriteString(aDiff.Text)
		}
	}
	return text.String()
}
