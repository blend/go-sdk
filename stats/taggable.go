/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

// Taggable is an interface for specifying and retrieving default stats tags
type Taggable interface {
	AddDefaultTag(string, string)
	AddDefaultTags(...string)
	DefaultTags() []string
}
