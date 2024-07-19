/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

// Selector is the common interface for selector types.
type Selector interface {
	Matches(labels Labels) bool
	Validate() error
	String() string
}
