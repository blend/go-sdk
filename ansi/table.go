package ansi

import (
	"fmt"
	"io"
	"reflect"
	"unicode/utf8"

	"github.com/blend/go-sdk/ex"
)

// Table character constants.
const (
	TableTopLeft     = '┌'
	TableTopRight    = '┐'
	TableBottomLeft  = '└'
	TableBottomRight = '┘'
	TableMidLeft     = '├'
	TableMidRight    = '┤'
	TableVertBar     = '│'
	TableHorizBar    = '─'
	TableTopSep      = '┬'
	TableBottomSep   = '┴'
	TableMidSep      = '┼'
)

// TableForSlice prints a table for a given slice.
// It will infer column names from the struct fields.
// If it is a mixed array (i.e. []interface{}) it will probably panic.
func TableForSlice(wr io.Writer, collection interface{}) error {
	// infer the column names from the fields
	cv := reflect.ValueOf(collection)
	for cv.Kind() == reflect.Ptr {
		cv = cv.Elem()
	}

	if cv.Kind() != reflect.Slice {
		return ex.New("table for slice; cannot iterate over non-slice collection")
	}

	ct := cv.Type()
	for ct.Kind() == reflect.Ptr || ct.Kind() == reflect.Slice {
		ct = ct.Elem()
	}

	columns := make([]string, ct.NumField())
	for index := 0; index < ct.NumField(); index++ {
		columns[index] = ct.Field(index).Name
	}

	var rows [][]string
	var rowValue reflect.Value
	for row := 0; row < cv.Len(); row++ {
		rowValue = cv.Index(row)
		rowValues := make([]string, ct.NumField())
		for fieldIndex := 0; fieldIndex < ct.NumField(); fieldIndex++ {
			rowValues[fieldIndex] = fmt.Sprintf("%v", rowValue.Field(fieldIndex).Interface())
		}
		rows = append(rows, rowValues)
	}

	return Table(wr, columns, rows)
}

// Table writes a table to a given writer.
func Table(wr io.Writer, columns []string, rows [][]string) error {
	if len(columns) == 0 {
		return ex.New("table; invalid columns; column set is empty")
	}

	/* begin establish max widths of columns */
	maxWidths := make([]int, len(columns))
	for index, columnName := range columns {
		maxWidths[index] = utf8.RuneCountInString(columnName)
	}
	for _, cols := range rows {
		for index, columnValue := range cols {
			if maxWidths[index] < utf8.RuneCountInString(columnValue) {
				maxWidths[index] = utf8.RuneCountInString(columnValue)
			}
		}
	}
	/* end establish max widths of columns */

	// draw headings
	io.WriteString(wr, string(TableTopLeft))
	for index := range columns {
		for x := 0; x < maxWidths[index]+2; x++ {
			io.WriteString(wr, string(TableHorizBar))
		}
		if index < (len(columns) - 1) {
			io.WriteString(wr, string(TableTopSep))
		}
	}
	io.WriteString(wr, string(TableTopRight))
	io.WriteString(wr, "\n")

	// draw column names
	io.WriteString(wr, string(TableVertBar))
	for index, column := range columns {
		writeWidth(wr, maxWidths[index], column)
		if index < (len(columns) - 1) {
			io.WriteString(wr, string(TableVertBar))
		}
	}
	io.WriteString(wr, string(TableVertBar))
	io.WriteString(wr, "\n")

	// draw bottom of column row
	io.WriteString(wr, string(TableMidLeft))
	for index := range columns {
		for x := 0; x < maxWidths[index]+2; x++ {
			io.WriteString(wr, string(TableHorizBar))
		}
		if index < (len(columns) - 1) {
			io.WriteString(wr, string(TableMidSep))
		}
	}
	io.WriteString(wr, string(TableMidRight))
	io.WriteString(wr, "\n")

	// draw rows
	for _, row := range rows {
		io.WriteString(wr, string(TableVertBar))
		for index, column := range row {
			writeWidth(wr, maxWidths[index], column)
			if index < (len(columns) - 1) {
				io.WriteString(wr, string(TableVertBar))
			}
		}
		io.WriteString(wr, string(TableVertBar))
		io.WriteString(wr, "\n")
	}

	// draw footer
	io.WriteString(wr, string(TableBottomLeft))
	for index := range columns {
		for x := 0; x < maxWidths[index]+2; x++ {
			io.WriteString(wr, string(TableHorizBar))
		}
		if index < (len(columns) - 1) {
			io.WriteString(wr, string(TableBottomSep))
		}
	}
	io.WriteString(wr, string(TableBottomRight))
	io.WriteString(wr, "\n")
	return nil
}

func writeWidth(wr io.Writer, width int, value string) (int, error) {
	return fmt.Fprintf(wr, fmt.Sprintf(" %%-%ds ", width), value)
}
