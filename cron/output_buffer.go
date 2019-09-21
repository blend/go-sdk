package cron

import (
	"bytes"
	"encoding/json"
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
	// Lines are the string lines broken up by newlines with associated timestamps
	Chunks []OutputChunk `json:"chunks"`
	// Listener is an optional listener for new line events.
	Listener OutputListener `json:"-"`
}

// Write writes the contents to the lines buffer.
func (ob *OutputBuffer) Write(contents []byte) (written int, err error) {
	chunk := OutputChunk{Timestamp: time.Now().UTC(), Data: contents}
	if ob.Listener != nil {
		ob.Listener(chunk)
	}
	ob.Chunks = append(ob.Chunks, chunk)
	written = len(contents)
	return
}

// Bytes rerturns the bytes written to the writer.
func (ob *OutputBuffer) Bytes() []byte {
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
