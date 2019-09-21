package cron

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"
)

// NewOutputBuffer creates a new line writer from a given set of bytes.
func NewOutputBuffer(contents []byte) *OutputBuffer {
	lw := new(OutputBuffer)
	lw.Write(contents)
	return lw
}

// OutputListener is a handler for lines.
type OutputListener func(OutputChunk)

// OutputChunk is a single write to the output buffer with a timestamp.
type OutputChunk struct {
	Timestamp time.Time
	Data      []byte
}

// Copy returns a copy of the chunk.
func (oc OutputChunk) Copy() OutputChunk {
	data := make([]byte, len(oc.Data))
	copy(data, oc.Data)
	return OutputChunk{
		Timestamp: oc.Timestamp,
		Data:      data,
	}
}

// MarshalJSON implemnts json.Marshaler.
func (oc OutputChunk) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"_ts":  oc.Timestamp,
		"data": string(oc.Data),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (oc *OutputChunk) UnmarshalJSON(contents []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(contents, &raw); err != nil {
		return err
	}

	if typed, ok := raw["_ts"].(string); ok {
		parsed, err := time.Parse(time.RFC3339, typed)
		if err != nil {
			return err
		}
		oc.Timestamp = parsed
	}
	if typed, ok := raw["data"].(string); ok {
		oc.Data = []byte(typed)
	}
	return nil
}

// OutputBuffer is a writer that accepts binary but splits out onto new lines.
type OutputBuffer struct {
	sync.RWMutex

	// Lines are the string lines broken up by newlines with associated timestamps
	Chunks []OutputChunk `json:"chunks"`
	// Listener is an optional listener for new line events.
	Listener OutputListener `json:"-"`
}

// Write writes the contents to the lines buffer.
// An important gotcha here is the `contents` parameter is by reference, as a result
// you can get into some bad loop capture states where this buffer will
// be assigned to multiple chunks but be effectively the same value.
// As a result, when you write to the output buffer, we fully copy this
// contents parameter for storage in the buffer.
func (ob *OutputBuffer) Write(contents []byte) (written int, err error) {
	data := make([]byte, len(contents))
	copy(data, contents)
	chunk := OutputChunk{Timestamp: time.Now().UTC(), Data: data}
	written = len(contents)

	// lock the buffer only to add the new chunk
	ob.Lock()
	ob.Chunks = append(ob.Chunks, chunk)
	ob.Unlock()

	// called outside critical section
	if ob.Listener != nil {
		// call the listener with a chunk copy.
		ob.Listener(chunk.Copy())
	}
	return
}

// Bytes rerturns the bytes written to the writer.
func (ob *OutputBuffer) Bytes() []byte {
	ob.RLock()
	defer ob.RUnlock()
	buffer := new(bytes.Buffer)
	for _, chunk := range ob.Chunks {
		buffer.Write(chunk.Data)
	}
	return buffer.Bytes()
}

// String returns the current combined output as a string.
func (ob *OutputBuffer) String() string {
	return string(ob.Bytes())
}
