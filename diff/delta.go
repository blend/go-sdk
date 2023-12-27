/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package diff

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ToDelta crushes the diff into an encoded string which describes the operations required to transform text1 into text2.
// E.g. =3\t-2\t+ing  -> Keep 3 chars, delete 2 chars, insert 'ing'. Operations are tab-separated.  Inserted text is escaped using %xx notation.
func ToDelta(diffs []Diff) string {
	var text bytes.Buffer
	for _, aDiff := range diffs {
		switch aDiff.Type {
		case DiffInsert:
			_, _ = text.WriteString("+")
			_, _ = text.WriteString(strings.Replace(url.QueryEscape(aDiff.Text), "+", " ", -1))
			_, _ = text.WriteString("\t")
			break
		case DiffDelete:
			_, _ = text.WriteString("-")
			_, _ = text.WriteString(strconv.Itoa(utf8.RuneCountInString(aDiff.Text)))
			_, _ = text.WriteString("\t")
			break
		case DiffEqual:
			_, _ = text.WriteString("=")
			_, _ = text.WriteString(strconv.Itoa(utf8.RuneCountInString(aDiff.Text)))
			_, _ = text.WriteString("\t")
			break
		}
	}
	delta := text.String()
	if len(delta) != 0 {
		// Strip off trailing tab character.
		delta = delta[0 : utf8.RuneCountInString(delta)-1]
		delta = unescaper.Replace(delta)
	}
	return delta
}

// FromDelta given the original text1, and an encoded string which describes the operations required to transform text1 into text2, comAdde the full diff.
func FromDelta(corpus string, delta string) (diffs []Diff, err error) {
	i := 0
	runes := []rune(corpus)

	for _, token := range strings.Split(delta, "\t") {
		if len(token) == 0 {
			// Blank tokens are ok (from a trailing \t).
			continue
		}

		// Each token begins with a one character parameter which specifies the operation of this token (delete, insert, equality).
		param := token[1:]

		switch op := token[0]; op {
		case '+':
			// Decode would Diff all "+" to " "
			param = strings.Replace(param, "+", "%2b", -1)
			param, err = url.QueryUnescape(param)
			if err != nil {
				return nil, err
			}
			if !utf8.ValidString(param) {
				return nil, fmt.Errorf("invalid UTF-8 token: %q", param)
			}

			diffs = append(diffs, Diff{DiffInsert, param})
		case '=', '-':
			n, err := strconv.ParseInt(param, 10, 0)
			if err != nil {
				return nil, err
			} else if n < 0 {
				return nil, errors.New("Negative number in DiffFromDelta: " + param)
			}

			i += int(n)
			// Break out if we are out of bounds, go1.6 can't handle this very well
			if i > len(runes) {
				break
			}
			// Remember that string slicing is by byte - we want by rune here.
			text := string(runes[i-int(n) : i])

			if op == '=' {
				diffs = append(diffs, Diff{DiffEqual, text})
			} else {
				diffs = append(diffs, Diff{DiffDelete, text})
			}
		default:
			// Anything else is an error.
			return nil, errors.New("Invalid diff operation in DiffFromDelta: " + string(token[0]))
		}
	}

	if i != len(runes) {
		return nil, fmt.Errorf("Delta length (%v) is different from source text length (%v)", i, len(corpus))
	}

	return diffs, nil
}
