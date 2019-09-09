package cron

import (
	"encoding/json"
	"time"
)

var (
	_ json.Marshaler   = (*OutputLine)(nil)
	_ json.Unmarshaler = (*OutputLine)(nil)
)

// FilterOutputLines applies a predicate to a set of lines.
func FilterOutputLines(lines []OutputLine, predicate func(OutputLine) bool) []OutputLine {
	var output []OutputLine
	for _, line := range lines {
		if predicate(line) {
			output = append(output, line)
		}
	}
	return output
}

// OutputLine is a line of output.
type OutputLine struct {
	Timestamp time.Time `json:"_ts"`
	Data      []byte    `json:"data"`
}

// MarshalJSON implements json.Marshaler.
func (l OutputLine) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"_ts":  l.Timestamp,
		"data": string(l.Data),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *OutputLine) UnmarshalJSON(contents []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(contents, &raw); err != nil {
		return err
	}

	if typed, ok := raw["_ts"].(string); ok {
		parsed, err := time.Parse(time.RFC3339, typed)
		if err != nil {
			return err
		}
		l.Timestamp = parsed
	}
	if typed, ok := raw["data"].(string); ok {
		l.Data = []byte(typed)
	}
	return nil
}
