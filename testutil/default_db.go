/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

import "github.com/blend/go-sdk/db"

var (
	_defaultDB *db.Connection
)

// DefaultDB returns a default database connection
// for tests.
func DefaultDB() *db.Connection {
	return _defaultDB
}
