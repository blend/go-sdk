package jobkit

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/blend/go-sdk/ex"
)

// NewLineWriter creates a new line writer from a given set of bytes.
func NewLineWriter(contents []byte) *LineWriter {
	lw := new(LineWriter)
	lw.Write(contents)
	return lw
}

// LineWriter is a writer that accepts binary but splits out onto new lines.
type LineWriter struct {
	Lines []Line

	// Current is a temporary holder for the current line.
	// It is added to the `Lines` slice when a newline is processed.
	Current Line
}

// Write writes the contents to the lines buffer.
func (lw *LineWriter) Write(contents []byte) (written int, err error) {
	var state int
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

// CombinedOutput returns the current combined output as a string.
func (lw *LineWriter) CombinedOutput() string {
	buffer := new(bytes.Buffer)
	for _, line := range lw.Lines {
		buffer.Write(line.Line)
		buffer.WriteRune('\n')
	}
	buffer.Write(lw.Current.Line)
	return buffer.String()
}

func (lw *LineWriter) commitLine() {
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
	return []byte(`{ "_ts": "` + l.Timestamp.Format(time.RFC3339) + `","line":"` + string(l.Line) + `"}`), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *Line) UnmarshalJSON(contents []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(contents, &raw); err != nil {
		return ex.New(err)
	}

	if typed, ok := raw["_ts"].(string); ok {
		parsed, err := time.Parse(time.RFC3339, typed)
		if err != nil {
			return ex.New(err)
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
