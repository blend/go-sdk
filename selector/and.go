/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import "strings"

// And is a combination selector.
type And []Selector

// Matches returns if both A and B match the labels.
func (a And) Matches(labels Labels) bool {
	for _, s := range a {
		if !s.Matches(labels) {
			return false
		}
	}
	return true
}

// Validate validates all the selectors in the clause.
func (a And) Validate() (err error) {
	for _, s := range a {
		err = s.Validate()
		if err != nil {
			return
		}
	}
	return
}

// And returns a string representation for the selector.
func (a And) String() string {
	var childValues []string
	for _, c := range a {
		childValues = append(childValues, c.String())
	}
	return strings.Join(childValues, ", ")
}
