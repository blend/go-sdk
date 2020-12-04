package fileutil

import (
	"crypto/rand"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFileReadByLines(t *testing.T) {
	assert := assert.New(t)

	f, err := NewTemp([]byte("this is a test\nof the emergency broadcast system\n"))
	assert.Nil(err)
	defer f.Close()

	var called, lineCorrect, readFirstLine bool
	err = ReadLines(f.Name(), func(line string) error {
		called = true
		if !readFirstLine {
			lineCorrect = line == "this is a test"
			readFirstLine = true
		}
		return nil
	})

	assert.Nil(err)
	assert.True(called, "We should have called the handler for `README.md`")
	assert.True(lineCorrect, "The first line should have matched the input")
}

func TestFileReadByChunks(t *testing.T) {
	assert := assert.New(t)

	buffer := make([]byte, 64)
	_, err := rand.Read(buffer)
	assert.Nil(err)

	f, err := NewTemp(buffer)
	assert.Nil(err)
	defer f.Close()

	var called, lengthCorrect bool
	err = ReadChunks(f.Name(), 32, func(chunk []byte) error {
		called = true
		lengthCorrect = len(chunk) == 32
		return nil
	})

	assert.Nil(err)
	assert.True(called, "We should have called the handler for `README.md`")
	assert.True(lengthCorrect, "We should have been passed 32 bytes")
}
