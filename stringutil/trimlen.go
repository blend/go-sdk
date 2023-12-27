/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

// TrimLen trims a string to a given length, i.e. the substring [0, length).
func TrimLen(val string, length int) string {
	if len(val) > length {
		return val[0:length]
	}
	return val
}
