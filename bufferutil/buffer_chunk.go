/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bufferutil

import (
	"encoding/json"
	"time"
)

// BufferChunkHandler is a handler for output chunks.
type BufferChunkHandler func(BufferChunk)

// BufferChunk is a single write to a buffer with a timestamp.
type BufferChunk struct {
	Timestamp	time.Time
	Data		[]byte
}

// Copy returns a copy of the chunk.
func (bc BufferChunk) Copy() BufferChunk {
	data := make([]byte, len(bc.Data))
	copy(data, bc.Data)
	return BufferChunk{
		Timestamp:	bc.Timestamp,
		Data:		data,
	}
}

// MarshalJSON implemnts json.Marshaler.
func (bc BufferChunk) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"_ts":	bc.Timestamp,
		"data":	string(bc.Data),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (bc *BufferChunk) UnmarshalJSON(contents []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(contents, &raw); err != nil {
		return err
	}

	if typed, ok := raw["_ts"].(string); ok {
		parsed, err := time.Parse(time.RFC3339, typed)
		if err != nil {
			return err
		}
		bc.Timestamp = parsed
	}
	if typed, ok := raw["data"].(string); ok {
		bc.Data = []byte(typed)
	}
	return nil
}
