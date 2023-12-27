/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"database/sql"

	"github.com/blend/go-sdk/ex"
)

// PoolCloseFinalizer is intended to be used in `defer` blocks with a named
// `error` return. It ensures a pool is closed after usage in contexts where
// a "limited use" pool is created.
//
// > func queries() (err error) {
// > 	var pool *db.Connection
// > 	defer func() {
// > 		err = db.PoolCloseFinalizer(pool, err)
// > 	}()
// > 	// ...
// > }
func PoolCloseFinalizer(pool *Connection, err error) error {
	if pool == nil || pool.Connection == nil {
		return err
	}

	closeErr := pool.Close()
	return ex.Nest(err, closeErr)
}

// TxRollbackFinalizer is intended to be used in `defer` blocks with a named
// `error` return. It ensures a transaction is always closed in blocks where
// a transaction is created.
//
// > func queries() (err error) {
// > 	var tx *sql.Tx
// > 	defer func() {
// > 		err = db.TxRollbackFinalizer(tx, err)
// > 	}()
// > 	// ...
// > }
func TxRollbackFinalizer(tx *sql.Tx, err error) error {
	if tx == nil {
		return err
	}

	rollbackErr := tx.Rollback()
	if rollbackErr == sql.ErrTxDone {
		return err
	}

	return ex.Nest(err, rollbackErr)
}
