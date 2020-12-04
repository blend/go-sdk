package stringutil

import (
	"math/rand"
	"time"
)

var (
	provider = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Random returns a random selection of runes from the set.
func Random(runeset []rune, length int) string {
	return Runeset(runeset).Random(length)
}
