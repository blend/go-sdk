/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package dbutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_ValidateDatabaseName(t *testing.T) {
	its := assert.New(t)

	testCases := [...]struct {
		Input	string
		Err	error
	}{
		{Input: "my_table"},
		{Input: "my_2nd_table"},
		{Input: "échéanciers"},
		{Input: "", Err: ErrDatabaseNameEmpty},
		{Input: strings.Repeat("a", DatabaseNameMaxLength+1), Err: ErrDatabaseNameTooLong},
		{Input: "2nd_table", Err: ErrDatabaseNameInvalidFirstRune},
		{Input: `"2nd_table"`, Err: ErrDatabaseNameInvalidFirstRune},
		{Input: "invalid-charater", Err: ErrDatabaseNameInvalid},
		{Input: "invalid'; DROP DB postgres; --", Err: ErrDatabaseNameInvalid},
		{Input: "postgres", Err: ErrDatabaseNameReserved},
		{Input: "template0", Err: ErrDatabaseNameReserved},
		{Input: "template1", Err: ErrDatabaseNameReserved},
		{Input: "defaultdb", Err: ErrDatabaseNameReserved},
	}

	var err error
	for _, tc := range testCases {
		err = ValidateDatabaseName(tc.Input)
		if tc.Err != nil && err != nil {
			its.Equal(tc.Err, ex.ErrClass(err))
		} else if tc.Err == nil && err != nil {
			its.FailNow(fmt.Sprintf("expected input %q not to produce an error, actual: %v", tc.Input, err))
		}
	}
}
