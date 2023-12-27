/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package dbutil

import "github.com/blend/go-sdk/db"

// OpenManagementConnection creates a database connection to the default database (typically postgres).
func OpenManagementConnection(options ...db.Option) (*db.Connection, error) {
	defaults := []db.Option{
		db.OptHost("localhost"),
		db.OptSSLMode("disable"),
		db.OptConfigFromEnv(),
		db.OptDatabase("postgres"),
	}
	conn, err := db.New(
		append(defaults, options...)...,
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
