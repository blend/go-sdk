/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

// GetDefaultsAsTags reads constants supplied at collector initialization and
// returns the keys and values formatted as a slice of stats tags.
func GetDefaultsAsTags(fixedTags map[string]string) (tags []string) {
	for key, value := range fixedTags {
		tags = append(tags, Tag(key, value))
	}
	return
}
