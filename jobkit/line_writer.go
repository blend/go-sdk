package jobkit

import (
	"time"
)

// LineWriter is a writer that accepts binary but splits out onto new lines.
type LineWriter struct {
	Buffer []byte
	Lines  []Line
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
			lw.Buffer = append(lw.Buffer, b)
		case 1:
			state = 0
			if b == '\n' {
				lw.commitLine()
				continue
			}
			lw.Buffer = append(lw.Buffer, b)
		}
	}
	written = len(contents)
	return
}

func (lw *LineWriter) commitLine() {
	lineBuffer := make([]byte, len(lw.Buffer))
	copy(lineBuffer, lw.Buffer)
	lw.Lines = append(lw.Lines, Line{
		Timestamp: time.Now().UTC(),
		Line:      string(lineBuffer),
	})
	lw.Buffer = nil
	return
}

// Line is a line of output.
type Line struct {
	Timestamp time.Time `json:"_ts"`
	Line      string    `json:"line"`
}
