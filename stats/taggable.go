/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stats

// Taggable is an interface for specifying and retrieving default stats tags
type Taggable interface {
	AddDefaultTag(string, string)
	AddDefaultTags(...string)
	DefaultTags() []string
}
