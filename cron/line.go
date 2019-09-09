package cron

import (
	"encoding/json"
	"time"
)

var (
	_ json.Marshaler   = (*Line)(nil)
	_ json.Unmarshaler = (*Line)(nil)
)

// FilterLines applies a predicate to a set of lines.
func FilterLines(lines []Line, predicate func(Line) bool) []Line {
	var output []Line
	for _, line := range lines {
		if predicate(line) {
			output = append(output, line)
		}
	}
	return output
}

// Line is a line of output.
type Line struct {
	Timestamp time.Time `json:"_ts"`
	Line      []byte    `json:"line"`
}

// MarshalJSON implements json.Marshaler.
func (l Line) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"_ts":  l.Timestamp,
		"line": string(l.Line),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *Line) UnmarshalJSON(contents []byte) error {
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
	if typed, ok := raw["line"].(string); ok {
		l.Line = []byte(typed)
	}
	return nil
}
