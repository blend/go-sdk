/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import "unicode/utf8"

// CheckKey validates a key.
func CheckKey(key string) (err error) {
	keyLen := len(key)
	if keyLen == 0 {
		err = ErrLabelKeyEmpty
		return
	}
	if keyLen > MaxLabelKeyTotalLen {
		err = ErrLabelKeyTooLong
		return
	}

	var working []rune
	var state int
	var ch rune
	var width int
	// separate the KEY into: DNS_SUBDOMAIN [ "/" DNS_LABEL ]
	for pos := 0; pos < keyLen; pos += width {
		ch, width = utf8.DecodeRuneInString(key[pos:])
		if state == 0 {
			if ch == ForwardSlash {
				err = CheckDNS(string(working))
				if err != nil {
					return
				}
				working = nil
				state = 1
				continue
			}
		}
		working = append(working, ch)
		continue
	}

	if len(working) == 0 {
		return ErrLabelKeyEmpty
	}
	if len(working) > MaxLabelKeyLen {
		return ErrLabelKeyTooLong
	}
	return CheckName(string(working))
}
