/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package diff

import "strconv"

// Operation defines the operation of a diff item.
type Operation int8

// Operation constants.
const (
	// DiffDelete item represents a delete diff.
	DiffDelete	Operation	= -1
	// DiffInsert item represents an insert diff.
	DiffInsert	Operation	= 1
	// DiffEqual item represents an equal diff.
	DiffEqual	Operation	= 0
	//IndexSeparator is used to separate the array indexes in an index string
	IndexSeparator	= ","
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DiffDelete - -1]
	_ = x[DiffInsert-1]
	_ = x[DiffEqual-0]
}

const operationName = "DeleteEqualInsert"

var operationIndex = [...]uint8{0, 6, 11, 17}

func (i Operation) String() string {
	i -= -1
	if i < 0 || i >= Operation(len(operationIndex)-1) {
		return "Operation(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return operationName[operationIndex[i]:operationIndex[i+1]]
}
