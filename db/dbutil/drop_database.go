/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package dbutil

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
)

// DropDatabase drops a database.
func DropDatabase(ctx context.Context, name string, opts ...db.Option) (err error) {
	var conn *db.Connection
	defer func() {
		err = db.PoolCloseFinalizer(conn, err)
	}()

	conn, err = OpenManagementConnection(opts...)
	if err != nil {
		return
	}

	if err = CloseAllConnections(ctx, conn, name); err != nil {
		return
	}

	_, err = conn.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s ", name))
	return
}
