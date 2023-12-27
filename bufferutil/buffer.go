/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package bufferutil

import (
	"bytes"
	"sync"
	"time"
)

// NewBuffer creates a new line writer from a given set of bytes.
func NewBuffer(contents []byte) *Buffer {
	lw := new(Buffer)
	_, _ = lw.Write(contents)
	return lw
}

// Buffer is a writer that accepts binary but splits out onto new lines.
type Buffer struct {
	sync.RWMutex
	// Size is the size of all bytes written to the buffer.
	Size	int64
	// Lines are the string lines broken up by newlines with associated timestamps
	Chunks	[]BufferChunk	`json:"chunks"`
	// Handler is an optional listener for new line events.
	Handler	BufferChunkHandler	`json:"-"`
}

// Write writes the contents to the output buffer.
// An important gotcha here is the `contents` parameter is by reference, as a result
// you can get into some bad loop capture states where this buffer will
// be assigned to multiple chunks but be effectively the same value.
// As a result, when you write to the output buffer, we fully copy this
// contents parameter for storage in the buffer.
func (b *Buffer) Write(contents []byte) (written int, err error) {
	chunkData := make([]byte, len(contents))
	copy(chunkData, contents)

	chunk := BufferChunk{Timestamp: time.Now().UTC(), Data: chunkData}
	written = len(chunkData)

	// lock the buffer only to add the new chunk
	b.Lock()
	b.Size += int64(written)
	b.Chunks = append(b.Chunks, chunk)
	b.Unlock()

	// called outside critical section
	if b.Handler != nil {
		// call the handler with the chunk.
		b.Handler(chunk)
	}
	return
}

// Bytes returns the bytes written to the writer.
func (b *Buffer) Bytes() []byte {
	b.RLock()
	defer b.RUnlock()
	buffer := new(bytes.Buffer)
	for _, chunk := range b.Chunks {
		buffer.Write(chunk.Data)
	}
	return buffer.Bytes()
}

// String returns the current combined output as a string.
func (b *Buffer) String() string {
	return string(b.Bytes())
}
