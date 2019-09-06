package stringutil

import (
	"bytes"
	"encoding/json"
	"time"
)

// NewLineBuffer creates a new line writer from a given set of bytes.
func NewLineBuffer(contents []byte) *LineBuffer {
	lw := new(LineBuffer)
	lw.Write(contents)
	return lw
}

// LineBuffer is a writer that accepts binary but splits out onto new lines.
type LineBuffer struct {
	// Lines are the string lines broken up by newlines with associated timestamps.
	Lines []Line
	// Current is a temporary holder for the current line.
	// It is added to the `Lines` slice when a newline is processed.
	Current Line
}

// Write writes the contents to the lines buffer.
func (lw *LineBuffer) Write(contents []byte) (written int, err error) {
	var state int
	if lw.Current.Timestamp.IsZero() {
		lw.Current.Timestamp = time.Now().UTC()
	}
	for _, b := range contents {
		switch state {
		case 0: //
			if b == '\r' { //handle CrLf
				state = 1
				continue
			}
			if b == '\n' {
				lw.commitLine()
				continue
			}
			lw.Current.Line = append(lw.Current.Line, b)
		case 1:
			state = 0
			if b == '\n' {
				lw.commitLine()
				continue
			}
			lw.Current.Line = append(lw.Current.Line, b)
		}
	}
	written = len(contents)
	return
}

// Bytes rerturns the bytes written to the writer.
func (lw *LineBuffer) Bytes() []byte {
	buffer := new(bytes.Buffer)
	for _, line := range lw.Lines {
		buffer.Write(line.Line)
		buffer.WriteRune('\n')
	}
	buffer.Write(lw.Current.Line)
	return buffer.Bytes()
}

// String returns the current combined output as a string.
func (lw *LineBuffer) String() string {
	return string(lw.Bytes())
}

// commit line adds the current line to the lines set and resets
// the current line
func (lw *LineBuffer) commitLine() {
	lw.Lines = append(lw.Lines, lw.Current)
	lw.Current.Timestamp = time.Time{}
	lw.Current.Line = nil
	return
}

var (
	_ json.Marshaler   = (*Line)(nil)
	_ json.Unmarshaler = (*Line)(nil)
)

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
