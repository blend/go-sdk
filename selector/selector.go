/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

// Selector is the common interface for selector types.
type Selector interface {
	Matches(labels Labels) bool
	Validate() error
	String() string
}
