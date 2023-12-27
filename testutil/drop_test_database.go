/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"

	"github.com/blend/go-sdk/db/dbutil"
)

// DropTestDatabase drops a database.
func DropTestDatabase(ctx context.Context, conn *db.Connection, opts ...db.Option) (err error) {
	var mgmt *db.Connection
	defer func() {
		err = db.PoolCloseFinalizer(mgmt, err)
	}()

	config, err := conn.Config.Reparse()
	if err != nil {
		return
	}

	mgmt, err = dbutil.OpenManagementConnection(opts...)
	if err != nil {
		return
	}

	_, err = mgmt.ExecContext(ctx, fmt.Sprintf("DROP DATABASE %s", config.Database))
	return
}
