package sh

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestMaxBytesWriter(t *testing.T) {
	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)
	mbw := LimitBytes(64, buf)

	written, err := mbw.Write(makeChunk(32))
	assert.Nil(err)
	assert.Equal(32, written)

	written, err = mbw.Write(makeChunk(16))
	assert.Nil(err)
	assert.Equal(16, written)

	written, err = mbw.Write(makeChunk(32))
	assert.True(ex.Is(err, ErrMaxBytesWriterCapacityLimit))
	assert.Equal(0, written)
}

func makeChunk(len int) []byte {
	output := make([]byte, len)
	for x := 0; x < len; x++ {
		output[x] = randomLetter()
	}
	return output
}

var (
	provider = rand.New(rand.NewSource(time.Now().UnixNano()))

	// LowerLetters is a runset of lowercase letters.
	lowerLetters = []byte("abcdefghijklmnopqrstuvwxyz")

	// UpperLetters is a runset of uppercase letters.
	upperLetters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Letters is a runset of both lower and uppercase letters.
	letters = append(lowerLetters, upperLetters...)
)

func randomLetter() byte {
	return letters[provider.Intn(len(letters))]
}
