/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
