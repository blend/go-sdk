/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"

	"github.com/blend/go-sdk/db"
)

const (
	dropTable     = "DROP TABLE IF EXISTS might_deadlock;"
	createTable   = "CREATE TABLE might_deadlock ( counter INTEGER NOT NULL, key TEXT NOT NULL );"
	tableSeedData = "INSERT INTO might_deadlock (counter, key) VALUES (4, 'hello'), (7, 'world'), (10, 'hello'), (5, 'world'), (3, 'world');"
)

func seedDatabase(ctx context.Context, pool *db.Connection) error {
	_, err := pool.ExecContext(ctx, dropTable)
	if err != nil {
		return err
	}

	_, err = pool.ExecContext(ctx, createTable)
	if err != nil {
		return err
	}

	_, err = pool.ExecContext(ctx, tableSeedData)
	return err
}
