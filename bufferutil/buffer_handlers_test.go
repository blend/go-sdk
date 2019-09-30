package bufferutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBufferHandlers(t *testing.T) {
	assert := assert.New(t)

	handlers := new(BufferHandlers)
	defer handlers.Close()

	datums := make(chan string, 2)

	didCallOne := make(chan struct{})
	handlers.Add("one", func(c BufferChunk) {
		datums <- string(c.Data)
		close(didCallOne)
	})

	didCallTwo := make(chan struct{})
	handlers.Add("two", func(c BufferChunk) {
		datums <- string(c.Data)
		close(didCallTwo)
	})

	go func() {
		handlers.Handle(BufferChunk{Data: []byte("hi")})
	}()

	<-didCallOne
	<-didCallTwo

	assert.Len(datums, 2)
}
