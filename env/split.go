/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
