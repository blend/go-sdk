package cron

import (
	"bytes"
	"time"
)

// NewLineBuffer creates a new line writer from a given set of bytes.
func NewLineBuffer(contents []byte) *LineBuffer {
	lw := new(LineBuffer)
	lw.Write(contents)
	return lw
}

// LineHandler is a handler for lines.
type LineHandler func(Line)

// LineBuffer is a writer that accepts binary but splits out onto new lines.
type LineBuffer struct {
	// Lines are the string lines broken up by newlines with associated timestamps.
	Lines []Line `json:"lines"`
	// Current is a temporary holder for the current line.
	// It is added to the `Lines` slice when a newline is processed.
	Current Line `json:"current"`
	// Handlers are listeners for new lines.
	Handlers *LineHandlers `json:"-"`
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
	lw.Handlers.Handle(lw.Current)

	lw.Lines = append(lw.Lines, lw.Current)
	lw.Current.Timestamp = time.Time{}
	lw.Current.Line = nil
	return
}
