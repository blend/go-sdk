package collections

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStringArray(t *testing.T) {
	a := assert.New(t)

	sa := Strings{"Foo", "bar", "baz"}

	a.Equal("Foo", sa.First())
	a.Equal("baz", sa.Last())

	a.True(sa.Contains("Foo"))
	a.False(sa.Contains("FOO"))
	a.False(sa.Contains("will"))

	a.True(sa.ContainsLower("foo"))
	a.False(sa.ContainsLower("will"))

	foo := sa.GetByLower("foo")
	a.Equal("Foo", foo)
	notFoo := sa.GetByLower("will")
	a.Equal("", notFoo)
}

func TestStringArrayReverse(t *testing.T) {
	a := assert.New(t)

	var rev Strings
	for arraySize := 0; arraySize < 13; arraySize++ {
		var arr Strings
		for x := 0; x < arraySize; x++ {
			arr = append(arr, strconv.Itoa(x))
		}
		rev = arr.Reverse()
		if arraySize == 0 {
			a.Empty(rev)
		} else if arraySize == 1 {
			a.Len(rev, 1)
			a.Equal(rev[0], arr[0])
		} else if arraySize > 1 {
			for y := 0; y < arraySize-1; y++ {
				a.Equal(rev[y], arr[arraySize-(y+1)], fmt.Sprintf("array size: %d", arraySize))
			}
		}
	}
}
