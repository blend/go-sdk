/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
