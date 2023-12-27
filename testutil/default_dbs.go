/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

import (
	"database/sql"

	"github.com/blend/go-sdk/db"
)

var (
	_defaultDBs []*db.Connection
)

// DefaultDBs returns a default set database connections for tests.
func DefaultDBs() []*db.Connection {
	return _defaultDBs
}

// BeginAll begins a transaction in each of the underlying connections.
// If an error is raised by *any* of the connections, t
func BeginAll() ([]*sql.Tx, error) {
	var output []*sql.Tx
	for x := 0; x < len(_defaultDBs); x++ {
		tx, err := _defaultDBs[x].Begin()
		if err != nil {
			for _, inFlight := range output {
				func() { _ = inFlight.Rollback() }()
			}
			return nil, err
		}
		output = append(output, tx)
	}
	return output, nil
}

// RollbackAll calls `Rollback` on a set of transactions.
func RollbackAll(txs ...*sql.Tx) error {
	for _, tx := range txs {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}
	return nil
}
