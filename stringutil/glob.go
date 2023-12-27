/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import "strings"

const (
	// GlobStar is the "match anything" constant.
	GlobStar = "*"
)

// GlobAny tests if a file matches a (potentially) csv of glob filters.
func GlobAny(subj string, patterns ...string) bool {
	for _, pattern := range patterns {
		if matches := Glob(subj, pattern); matches {
			return true
		}
	}
	return false
}

// Glob returns if a subject matches a given pattern.
func Glob(subj, pattern string) bool {
	if pattern == "" {
		return subj == pattern
	}

	if pattern == GlobStar {
		return true
	}

	parts := strings.Split(pattern, GlobStar)

	if len(parts) == 1 {
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, GlobStar)
	trailingGlob := strings.HasSuffix(pattern, GlobStar)
	end := len(parts) - 1

	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			if idx < 0 {
				return false
			}
		}

		subj = subj[idx+len(parts[i]):]
	}

	return trailingGlob || strings.HasSuffix(subj, parts[end])
}
