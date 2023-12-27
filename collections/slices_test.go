/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package collections

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFirst(t *testing.T) {
	a := assert.New(t)
	sa := []string{"Foo", "bar", "baz"}
	a.Equal("Foo", First(sa))
}

func TestLast(t *testing.T) {
	a := assert.New(t)
	sa := []string{"Foo", "bar", "baz"}
	a.Equal("baz", Last(sa))
}

func TestContains(t *testing.T) {
	a := assert.New(t)
	sa := []string{"Foo", "bar", "baz"}
	a.True(Contains(sa, "Foo"))
	a.False(Contains(sa, "FOO"))
	a.False(Contains(sa, "will"))
}

func TestReverse(t *testing.T) {
	a := assert.New(t)
	strSlice := []string{"h", "e", "l", "l", "o"}
	strSlice2 := Reverse(strSlice)
	a.Equal([]string{"h", "e", "l", "l", "o"}, strSlice)
	a.Equal([]string{"o", "l", "l", "e", "h"}, strSlice2)

	intSlice := []int{}
	newSlice := Reverse(intSlice)
	a.Nil(newSlice)

	intSlice = []int{0}
	newSlice = Reverse(intSlice)
	a.Equal([]int{0}, newSlice)

	intSlice = []int{0, 1}
	newSlice = Reverse(intSlice)
	a.Equal([]int{0, 1}, intSlice)
	a.Equal([]int{1, 0}, newSlice)

	intSlice = []int{0, 1, 2}
	newSlice = Reverse(intSlice)
	a.Equal([]int{0, 1, 2}, intSlice)
	a.Equal([]int{2, 1, 0}, newSlice)

	intSlice = []int{0, 1, 2, 3}
	newSlice = Reverse(intSlice)
	a.Equal([]int{0, 1, 2, 3}, intSlice)
	a.Equal([]int{3, 2, 1, 0}, newSlice)

	intSlice = []int{0, 1, 2, 3, 4}
	newSlice = Reverse(intSlice)
	a.Equal([]int{0, 1, 2, 3, 4}, intSlice)
	a.Equal([]int{4, 3, 2, 1, 0}, newSlice)
}

func TestReverseInPlace(t *testing.T) {
	a := assert.New(t)
	strSlice := []string{"h", "e", "l", "l", "o"}
	ReverseInPlace(strSlice)
	a.Equal([]string{"o", "l", "l", "e", "h"}, strSlice)

	intSlice := []int{}
	ReverseInPlace(intSlice)
	a.Equal([]int{}, intSlice)

	intSlice = []int{0}
	ReverseInPlace(intSlice)
	a.Equal([]int{0}, intSlice)

	intSlice = []int{0, 1}
	ReverseInPlace(intSlice)
	a.Equal([]int{1, 0}, intSlice)

	intSlice = []int{0, 1, 2}
	ReverseInPlace(intSlice)
	a.Equal([]int{2, 1, 0}, intSlice)

	intSlice = []int{0, 1, 2, 3}
	ReverseInPlace(intSlice)
	a.Equal([]int{3, 2, 1, 0}, intSlice)

	intSlice = []int{0, 1, 2, 3, 4}
	ReverseInPlace(intSlice)
	a.Equal([]int{4, 3, 2, 1, 0}, intSlice)
}
