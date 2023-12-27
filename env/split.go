/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env

// Split splits an env var in the form `KEY=value`.
func Split(s string) (key, value string) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			key = s[:i]
			value = s[i+1:]
			return
		}
	}
	return
}
