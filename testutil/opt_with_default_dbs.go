/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"context"

	"github.com/blend/go-sdk/db"
)

// OptWithDefaultDBs runs a test suite with a count of database connections.
// Note: this type of connection pool is used in rare circumstances for
// performance reasons; you probably want to use `OptWithDefaultDB` for your tests.
func OptWithDefaultDBs(count int) Option {
	return func(s *Suite) {
		s.Before = append(s.Before, func(ctx context.Context) error {
			_defaultDBs = make([]*db.Connection, count)
			for index := 0; index < count; index++ {
				conn, err := CreateTestDatabase(ctx)
				if err != nil {
					return err
				}
				_defaultDBs[index] = conn
			}
			return nil
		})
		s.After = append(s.After, func(ctx context.Context) error {
			for index := range _defaultDBs {
				if err := _defaultDBs[index].Close(); err != nil {
					return err
				}
				if err := DropTestDatabase(ctx, _defaultDBs[index]); err != nil {
					return err
				}
			}
			return nil
		})
	}
}
