/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ansi

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/blend/go-sdk/ex"
)

// Table character constants.
const (
	TableTopLeft		= "┌"
	TableTopRight		= "┐"
	TableBottomLeft		= "└"
	TableBottomRight	= "┘"
	TableMidLeft		= "├"
	TableMidRight		= "┤"
	TableVertBar		= "│"
	TableHorizBar		= "─"
	TableTopSep		= "┬"
	TableBottomSep		= "┴"
	TableMidSep		= "┼"
	NewLine			= "\n"
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
func Table(wr io.Writer, columns []string, rows [][]string) (err error) {
	if len(columns) == 0 {
		return ex.New("table; invalid columns; column set is empty")
	}

	/* helpers */
	defer func() {
		if r := recover(); r != nil {
			if typed, ok := r.(error); ok {
				err = typed
			}
		}
	}()
	write := func(str string) {
		_, writeErr := io.WriteString(wr, str)
		if writeErr != nil {
			panic(writeErr)
		}
	}
	/* end helpers */

	/* begin establish max widths of columns */
	maxWidths := make([]int, len(columns))
	for index, columnName := range columns {
		maxWidths[index] = stringWidth(columnName)
	}

	var width int
	for _, cols := range rows {
		for index, columnValue := range cols {
			width = stringWidth(columnValue)
			if maxWidths[index] < width {
				maxWidths[index] = width
			}
		}
	}
	/* end establish max widths of columns */

	/* draw top of column row */
	write(TableTopLeft)
	for index := range columns {
		write(repeat(TableHorizBar, maxWidths[index]))
		if isNotLast(index, columns) {
			write(TableTopSep)
		}
	}
	write(TableTopRight)
	write(NewLine)
	/* end draw top of column row */

	/* draw column names */
	write(TableVertBar)
	for index, columnLabel := range columns {
		write(padRight(columnLabel, maxWidths[index]))
		if isNotLast(index, columns) {
			write(TableVertBar)
		}
	}
	write(TableVertBar)
	write(NewLine)
	/* end draw column names */

	/* draw bottom of column row */
	write(TableMidLeft)
	for index := range columns {
		write(repeat(TableHorizBar, maxWidths[index]))
		if isNotLast(index, columns) {
			write(TableMidSep)
		}
	}
	write(TableMidRight)
	write(NewLine)
	/* end draw bottom of column row */

	/* draw rows */
	for _, row := range rows {
		write(TableVertBar)
		for index, column := range row {
			write(padRight(column, maxWidths[index]))
			if isNotLast(index, columns) {
				write(TableVertBar)
			}
		}
		write(TableVertBar)
		write(NewLine)
	}
	/* end draw rows */

	/* draw footer */
	write(TableBottomLeft)
	for index := range columns {
		write(repeat(TableHorizBar, maxWidths[index]))
		if isNotLast(index, columns) {
			write(TableBottomSep)
		}
	}
	write(TableBottomRight)
	write(NewLine)
	/* end draw footer */
	return
}

func stringWidth(value string) (width int) {
	var runeWidth int
	for _, c := range value {
		runeWidth = utf8.RuneLen(c)
		if runeWidth > 1 {
			width += 2
		} else {
			width++
		}
	}
	return
}

func repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

func padRight(value string, width int) string {
	valueWidth := stringWidth(value)
	spaces := width - valueWidth
	if spaces == 0 {
		return value
	}
	return value + strings.Repeat(" ", spaces)
}

func isNotLast(index int, values []string) bool {
	return index < (len(values) - 1)
}
