/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSet(t *testing.T) {
	a := assert.New(t)

	set := Set[int]{}
	set.Add(1)
	a.True(set.Contains(1))
	a.Equal(1, set.Len())
	a.False(set.Contains(2))
	set.Remove(1)
	a.False(set.Contains(1))
	a.Zero(set.Len())
}

func TestSetOperations(t *testing.T) {
	a := assert.New(t)

	s1 := NewSet[int](1, 2, 3, 4)
	s2 := NewSet[int](1, 2)
	s3 := NewSet[int](3, 4, 5, 6)

	union := s1.Union(s3)
	a.Len(union, 6)
	intersect := s1.Intersect(s2)
	a.Len(intersect, 2)
	diff := s1.Difference(s3)
	a.Len(diff, 4)
	diff = s3.Difference(s1)
	a.Len(diff, 4)
	a.True(s2.IsSubsetOf(s1))
	a.False(s1.IsSubsetOf(s2))
}
