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
	dropTable     = "DROP TABLE IF EXISTS might_sleep;"
	createTable   = "CREATE TABLE might_sleep ( id INTEGER NOT NULL );"
	tableSeedData = "INSERT INTO might_sleep (id) VALUES (1337);"
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
