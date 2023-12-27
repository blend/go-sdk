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
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/dbutil"
	"github.com/blend/go-sdk/uuid"
)

// CreateTestDatabase creates a randomized test database.
func CreateTestDatabase(ctx context.Context, opts ...db.Option) (*db.Connection, error) {
	databaseName := fmt.Sprintf("testdb_%s", uuid.V4().String())
	if err := dbutil.CreateDatabase(ctx, databaseName, opts...); err != nil {
		return nil, err
	}

	defaults := []db.Option{
		db.OptHost("localhost"),
		db.OptSSLMode("disable"),
		db.OptConfigFromEnv(),
		db.OptDatabase(databaseName),
		db.OptDialect(db.DialectPostgres),
	}
	conn, err := db.New(
		append(defaults, opts...)...,
	)
	if err != nil {
		return nil, err
	}
	err = conn.Open()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
