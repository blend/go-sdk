package cron

import (
	"bytes"
	"time"
)

// NewOutputBuffer creates a new line writer from a given set of bytes.
func NewOutputBuffer(contents []byte) *OutputBuffer {
	lw := new(OutputBuffer)
	lw.Write(contents)
	return lw
}

// OutputListener is a handler for lines.
type OutputListener func(OutputLine)

// OutputBuffer is a writer that accepts binary but splits out onto new lines.
type OutputBuffer struct {
	// Lines are the string lines broken up by newlines with associated timestamps
	Lines []OutputLine `json:"lines"`
	// Working is a temporary holder for the current line.
	// It is added to the `Lines` slice when a newline is processed.
	Working OutputLine `json:"working"`
	// Listener is an optional listener for new line events.
	Listener OutputListener `json:"-"`
}

// Write writes the contents to the lines buffer.
func (lw *OutputBuffer) Write(contents []byte) (written int, err error) {
	var state int
	if lw.Working.Timestamp.IsZero() {
		lw.Working.Timestamp = time.Now().UTC()
	}
	for _, b := range contents {
		switch state {
		case 0: //
			if b == '\r' { //handle CrLf
				state = 1
				continue
			}
			if b == '\n' {
				lw.AddWorkingLine()
				continue
			}
			lw.Working.Data = append(lw.Working.Data, b)
		case 1:
			state = 0
			if b == '\n' {
				lw.AddWorkingLine()
				continue
			}
			lw.Working.Data = append(lw.Working.Data, b)
		}
	}
	written = len(contents)
	return
}

// Bytes rerturns the bytes written to the writer.
func (lw *OutputBuffer) Bytes() []byte {
	buffer := new(bytes.Buffer)
	for _, line := range lw.Lines {
		buffer.Write(line.Data)
		buffer.WriteRune('\n')
	}
	buffer.Write(lw.Working.Data)
	return buffer.Bytes()
}

// String returns the current combined output as a string.
func (lw *OutputBuffer) String() string {
	return string(lw.Bytes())
}

// AddWorkingLine adds the current line to the lines set and resets
// the current line
func (lw *OutputBuffer) AddWorkingLine() {
	if lw.Listener != nil {
		lw.Listener(lw.Working)
	}
	lw.Lines = append(lw.Lines, lw.Working)
	lw.Working.Timestamp = time.Time{}
	lw.Working.Data = nil
	return
}
