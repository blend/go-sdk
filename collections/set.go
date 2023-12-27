/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

import (
	"fmt"
	"strings"
)

// NewSet creates a new set from a list of values.
func NewSet[T comparable](values ...T) Set[T] {
	output := Set[T](make(map[T]struct{}))
	for _, v := range values {
		output[v] = struct{}{}
	}
	return output
}

// Set is a generic implementation of a set data structure.
type Set[T comparable] map[T]struct{}

// Add adds an element to the set, replacing a previous value.
func (s Set[T]) Add(i T) {
	s[i] = struct{}{}
}

// Remove removes an element from the set.
func (s Set[T]) Remove(i T) {
	delete(s, i)
}

// Contains returns if the element is in the set.
func (s Set[T]) Contains(i T) bool {
	_, ok := s[i]
	return ok
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// Copy returns a new copy of the set.
func (s Set[T]) Copy() Set[T] {
	newSet := NewSet[T]()
	for key := range s {
		newSet.Add(key)
	}
	return newSet
}

// Union joins two sets together without dupes.
func (s Set[T]) Union(other Set[T]) Set[T] {
	union := NewSet[T]()
	for k := range s {
		union.Add(k)
	}

	for k := range other {
		union.Add(k)
	}
	return union
}

// Intersect returns shared elements between two sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	intersection := NewSet[T]()
	for k := range s {
		if other.Contains(k) {
			intersection.Add(k)
		}
	}
	return intersection
}

// Difference returns non-shared elements between two sets.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	difference := NewSet[T]()
	for k := range s {
		if !other.Contains(k) {
			difference.Add(k)
		}
	}
	for k := range other {
		if !s.Contains(k) {
			difference.Add(k)
		}
	}
	return difference
}

// IsSubsetOf returns if a given set is a complete subset of another set,
// i.e. all elements in target set are in other set.
func (s Set[T]) IsSubsetOf(other Set[T]) bool {
	for k := range s {
		if !other.Contains(k) {
			return false
		}
	}
	return true
}

// AsSlice returns the set as a slice.
func (s Set[T]) AsSlice() []T {
	var output []T
	for key := range s {
		output = append(output, key)
	}
	return output
}

// String returns the set as a csv string.
func (s Set[T]) String() string {
	var values []string
	for key := range s {
		values = append(values, fmt.Sprint(key))
	}
	return strings.Join(values, ", ")
}
