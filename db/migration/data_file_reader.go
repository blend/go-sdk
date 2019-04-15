package migration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

const (
	runeTab          = rune('\t')
	runeNewline      = rune('\n')
	byteNewLine      = byte('\n')
	byteTab          = byte('\t')
	null             = `\N`
	regexCopyExtract = `COPY (.*)? \((.*)?\)`
)

// ReadDataFile returns a new DataFileReader
func ReadDataFile(filePath string) *DataFileReader {
	return &DataFileReader{
		path:          filePath,
		copyExtractor: regexp.MustCompile(regexCopyExtract),
	}
}

// DataFileReader reads a postgres dump.
type DataFileReader struct {
	path          string
	copyExtractor *regexp.Regexp
}

// Label returns the label for the data file reader.
func (dfr *DataFileReader) Label() string {
	return fmt.Sprintf("read data file `%s`", dfr.path)
}

// Action applies the file reader.
func (dfr *DataFileReader) Action(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
	var f *os.File
	if f, err = os.Open(dfr.path); err != nil {
		return
	}
	defer f.Close()

	var stmt *sql.Stmt
	var state int

	var cursor int64
	var readBuffer = make([]byte, 32)
	var readErr error
	var lineBuffer = bytes.NewBuffer([]byte{})
	var line string
	var pieces []interface{}

	for readErr == nil {
		lineBuffer.Reset()

		switch state {
		case 0:
			cursor, readErr = dfr.readLine(f, cursor, readBuffer, lineBuffer)
			if readErr != nil {
				continue
			}
			line = lineBuffer.String()

			if stringutil.HasPrefixCaseless(line, "--") {
				continue
			}

			if stringutil.HasPrefixCaseless(line, "set") {
				continue
			}

			if stringutil.HasPrefixCaseless(line, "copy") {
				if !stringutil.HasSuffixCaseless(line, "from stdin;") {
					err = fmt.Errorf("only `stdin` from clauses supported at this time, cannot continue")
					return
				}

				stmt, err = dfr.executeCopyLine(line, c, tx)
				if err != nil {
					return
				}
				state = 1
				continue
			}

			err = c.Invoke(db.OptTx(tx)).Exec(line)
			if err != nil {
				return
			}
		case 1:
			pieces, cursor, readErr = dfr.readTabLine(f, cursor, readBuffer, lineBuffer)
			if readErr != nil {
				continue
			}

			if len(pieces) == 0 {
				err = fmt.Errorf("empty data line, we error on this for now")
				return
			}

			if len(pieces) == 1 && stringutil.HasPrefixCaseless(pieces[0].(string), `\.`) {
				err = stmt.Close()
				if err != nil {
					return
				}
				state = 0
				continue
			}

			_, err = stmt.Exec(pieces...)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (dfr *DataFileReader) executeCopyLine(line string, c *db.Connection, tx *sql.Tx) (*sql.Stmt, error) {
	pieces := dfr.extractCopyLine(line)
	if len(pieces) < 3 {
		return nil, ex.New("Invalid `COPY ...` line, cannot continue.")
	}
	tableName := pieces[1]
	columnCSV := pieces[2]
	columns := strings.Split(columnCSV, ", ")
	return tx.Prepare(CopyIn(tableName, columns...))
}

// regexExtractSubMatches returns sub matches for an expr because go's regexp library is weird.
func (dfr *DataFileReader) extractCopyLine(line string) []string {
	allResults := dfr.copyExtractor.FindAllStringSubmatch(line, -1)
	results := []string{}
	for _, resultSet := range allResults {
		for _, result := range resultSet {
			results = append(results, result)
		}
	}

	return results
}

func (dfr *DataFileReader) extractDataLine(line string) []interface{} {
	var values []interface{}
	var value string
	var state int

	appendValue := func() {
		if value == `\N` {
			values = append(values, nil)
		} else {
			values = append(values, value)
		}
	}

	for _, r := range line {
		switch state {
		case 0:
			if r == runeTab {
				continue
			}
			state = 1
			value = value + string(r)
		case 1:
			if r == runeTab {
				appendValue()
				state = 0
				value = ""
				continue
			}

			value = value + string(r)
		}
	}

	if len(value) > 0 {
		appendValue()
	}

	return values
}

// readLine reads a file until a newline.
func (dfr *DataFileReader) readLine(f io.ReaderAt, cursor int64, readBuffer []byte, lineBuffer *bytes.Buffer) (int64, error) {
	// bytesRead is the return from the ReadAt function
	// it indicates how many effective bytes we read from the stream.
	var bytesRead int
	// err is our primary indicator if there was an issue with the stream
	// or if we've reached the end of the file.
	var err error
	// b is the byte we're reading at a time.
	var b byte

	// while we haven't hit an error (this includes EOF!)
	for err == nil {
		// read the stream
		bytesRead, err = f.ReadAt(readBuffer, cursor)
		// abort on error
		if err != nil && err != io.EOF { //let this continue on eof
			return cursor, err
		}

		// slurp the read buffer.
		for readBufferIndex := 0; readBufferIndex < bytesRead; readBufferIndex++ {
			// advance the cursor regardless of what we read out.
			// if we read a newline, great! we'll start the next character after the newline after.
			cursor++

			// slurp the byte out of the read buffer
			b = readBuffer[readBufferIndex]
			if b == byteNewLine {
				// we bifurcate here because we need to forward the eof
				// if we read the buffer exactly right.
				if readBufferIndex == bytesRead-1 {
					return cursor, err
				}
				// otherwise the newline may have happened
				// before the actual eof.
				return cursor, nil
			}

			// b wasnt a newline, write it to the output buffer.
			lineBuffer.WriteByte(b)
		}
	}
	// we've reached the end of the file
	// there may not have been a newline
	// return what we have
	return cursor, err
}

// readTabLine reads a file until a new line, collecting tab delimited sections into an array.
func (dfr *DataFileReader) readTabLine(f io.ReaderAt, cursor int64, readBuffer []byte, lineBuffer *bytes.Buffer) ([]interface{}, int64, error) {
	// bytesRead is the return from the ReadAt function
	// it indicates how many effective bytes we read from the stream.
	var bytesRead int
	// err is our primary indicator if there was an issue with the stream
	// or if we've reached the end of the file.
	var err error
	// b is the byte we're reading at a time.
	var b byte
	// pieces is used to collect the tab delimited components of the line.
	var pieces []interface{}
	// while we haven't hit an error (this includes EOF!)
	for err == nil {
		// read the stream
		bytesRead, err = f.ReadAt(readBuffer, cursor)
		// abort on error
		if err != nil && err != io.EOF { //let this continue on eof
			return pieces, cursor, err
		}

		// slurp the read buffer.
		for readBufferIndex := 0; readBufferIndex < bytesRead; readBufferIndex++ {
			// advance the cursor regardless of what we read out.
			// if we read a newline, great! we'll start the next character after the newline after.
			cursor++

			// slurp the byte out of the read buffer
			b = readBuffer[readBufferIndex]
			if b == byteNewLine {
				// make sure to collect the remaining text in the
				// linebuffer.
				pieces = dfr.readTabLineAppendPiece(pieces, lineBuffer)

				// we bifurcate here because we need to forward the eof
				// if we read the buffer exactly right.
				if readBufferIndex == bytesRead-1 {
					return pieces, cursor, err
				}
				// otherwise the newline may have happened
				// before the actual eof.
				return pieces, cursor, nil
			}

			// if we see a tab
			// mark a tab section
			// reset the buffer
			if b == byteTab {
				pieces = dfr.readTabLineAppendPiece(pieces, lineBuffer)
				continue
			}

			// b wasnt a newline, write it to the output buffer.
			lineBuffer.WriteByte(b)
		}
	}

	pieces = dfr.readTabLineAppendPiece(pieces, lineBuffer)
	// we've reached the end of the file
	// there may not have been a newline
	// return what we have
	return pieces, cursor, err
}

// readTabLineAppendPiece is a commonly used code block
// that conditionally adds a new piece to the tab piece collection
func (dfr *DataFileReader) readTabLineAppendPiece(pieces []interface{}, lineBuffer *bytes.Buffer) []interface{} {
	if lineBuffer.Len() > 0 {
		value := lineBuffer.String()
		lineBuffer.Reset()
		if value == null {
			return append(pieces, nil)
		}
		return append(pieces, value)
	}
	return pieces
}

// CopyIn creates a COPY FROM statement which can be prepared with
// Tx.Prepare().  The target table should be visible in search_path.
func CopyIn(table string, columns ...string) string {
	stmt := "COPY " + QuoteIdentifier(table) + " ("
	for i, col := range columns {
		if i != 0 {
			stmt += ", "
		}
		stmt += QuoteIdentifier(col)
	}
	stmt += ") FROM STDIN"
	return stmt
}

// QuoteIdentifier quotes an "identifier" (e.g. a table or a column name) to be
// used as part of an SQL statement.  For example:
//
//    tblname := "my_table"
//    data := "my_data"
//    quoted := pq.QuoteIdentifier(tblname)
//    err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUES ($1)", quoted), data)
//
// Any double quotes in name will be escaped.  The quoted identifier will be
// case sensitive when used in a query.  If the input string contains a zero
// byte, the result will be truncated immediately before it.
func QuoteIdentifier(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	return `"` + strings.Replace(name, `"`, `""`, -1) + `"`
}
