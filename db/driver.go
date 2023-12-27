/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

/*
USAGE NOTE from: https://github.com/jackc/pgx/blob/master/README.md#choosing-between-the-pgx-and-databasesql-interfaces

The database/sql interface only allows the underlying driver to return or
receive the following types: int64, float64, bool, []byte, string, time.Time, or nil.
Handling other types requires implementing the database/sql.Scanner and the
database/sql/driver/driver.Valuer interfaces which require transmission of values in text format.

The binary format can be substantially faster, which is what the pgx interface uses.
*/

import (
	// the default driver is the stdlib version of pgx
	_ "github.com/jackc/pgx/v4/stdlib"
)
