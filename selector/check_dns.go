/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import "unicode/utf8"

// CheckDNS returns if a given value is a conformant DNS_SUBDOMAIN.
//
// See: https://www.ietf.org/rfc/rfc952.txt and https://www.ietf.org/rfc/rfc1123.txt (2.1) for more information.
//
// Specifically a DNS_SUBDOMAIN must:
// - Be constituted of ([a-z0-9\-\.])
// - It must start and end with ([a-z0-9])
// - it must be less than 254 characters in length
// - Characters ('.', '-') cannot repeat
func CheckDNS(value string) (err error) {
	valueLen := len(value)
	if valueLen == 0 {
		err = ErrLabelKeyDNSSubdomainEmpty
		return
	}
	if valueLen > MaxLabelKeyDNSSubdomainLen {
		err = ErrLabelKeyDNSSubdomainTooLong
		return
	}

	var state int
	var ch rune
	var width int

	const (
		statePrefixSuffix = 0
		stateAlpha        = 1
		stateDotDash      = 2
	)

	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case statePrefixSuffix:
			if !isDNSAlpha(ch) {
				return ErrLabelKeyInvalidDNSSubdomain
			}
			state = stateAlpha
			continue

		// the last character was a ...
		case stateAlpha:
			if ch == Dot || ch == Dash {
				state = stateDotDash
				continue
			}
			if !isDNSAlpha(ch) {
				err = ErrLabelKeyInvalidDNSSubdomain
				return
			}
			// if we're at the penultimate char
			if pos == valueLen-2 {
				state = statePrefixSuffix
				continue
			}
			continue

		// the last character was a ...
		case stateDotDash:
			if isDNSAlpha(ch) {
				state = stateAlpha
				continue
			}
			err = ErrLabelKeyInvalidDNSSubdomain
			return
		}
	}
	return nil
}
