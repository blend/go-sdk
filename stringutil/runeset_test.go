package stringutil

import (
	"sort"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRunesetSort(t *testing.T) {
	assert := assert.New(t)

	sorted := Runeset([]rune("fedcba"))
	sort.Sort(sorted)
	assert.Equal([]rune("abcdef"), sorted)
}

func TestRunesetCombine(t *testing.T) {
	assert := assert.New(t)

	combined := Letters.Combine(Numbers, Symbols, Letters)
	assert.Len(combined, 84)
}

func TestRunesetRandom(t *testing.T) {
	assert := assert.New(t)

	output := LettersAndNumbers.Random(32)
	assert.Len(output, 32)
}
