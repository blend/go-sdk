/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import "bytes"

// Tokens is a soft alias to map[string]string
type Tokens = map[string]string

// Tokenize replaces a given set of tokens in a corpus.
// Tokens should appear in the corpus in the form ${[KEY]} where [KEY] is the key in the map.
// Examples: corpus: "foo/${bar}/baz", { "bar": "example-string" } => "foo/example-string/baz"
// UTF-8 is handled via. runes.
func Tokenize(corpus string, tokens Tokens) string {
	// there is no way to escape anything smaller than [3] b/c len("${}") == 3
	if len(corpus) < 3 {
		return corpus
	}
	// sanity check on tokens collection.
	if len(tokens) == 0 {
		return corpus
	}

	output := bytes.NewBuffer(nil)

	start0 := rune('$')
	start1 := rune('{')
	end0 := rune('}')

	var state int
	// working token is the full token (including ${ and }).
	// working key is the stuff within the ${ and }.
	var workingToken, workingKey *bytes.Buffer
	var key string

	for _, c := range corpus {
		switch state {
		case 0:	// non-token, add to output
			if c == start0 {
				state = 1
				workingToken = bytes.NewBuffer(nil)
				workingToken.WriteRune(c)
				continue
			}
			output.WriteRune(c)
			continue
		case 1:
			if c == start1 {
				state = 2	//consume token key
				workingToken.WriteRune(c)
				workingKey = bytes.NewBuffer(nil)
				continue
			}
			state = 0
			output.WriteString(workingToken.String())
			output.WriteRune(c)
			workingToken = nil
			workingKey = nil
			continue
		case 2:
			if c == end0 {
				workingToken.WriteRune(c)
				// lookup replacement token.
				key = workingKey.String()
				if value, hasValue := tokens[key]; hasValue {
					output.WriteString(value)
				} else {
					output.WriteString(workingToken.String())
				}
				workingToken = nil
				workingKey = nil
				state = 0
				continue
			}
			if c == start0 {
				state = 3
				workingToken.WriteRune(c)
				workingKey.WriteRune(c)
				continue
			}
			workingToken.WriteRune(c)
			workingKey.WriteRune(c)
			continue
		case 3:
			if c == start1 {
				state = 4
				workingToken.WriteRune(c)
				workingKey.WriteRune(c)
				continue
			}
			state = 2
			workingToken.WriteRune(c)
			workingKey.WriteRune(c)
			continue
		case 4:
			if c == end0 {
				state = 2
				workingToken.WriteRune(c)
				workingKey.WriteRune(c)
				continue
			}
			workingToken.WriteRune(c)
			workingKey.WriteRune(c)
			continue
		}
	}

	return output.String()
}
