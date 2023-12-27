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
)

// OptWithDefaultDB runs a test suite with a dedicated database connection.
func OptWithDefaultDB() Option {
	return func(s *Suite) {
		var err error
		s.Before = append(s.Before, func(ctx context.Context) error {
			_defaultDB, err = CreateTestDatabase(ctx)
			if err != nil {
				return err
			}
			return nil
		})
		s.After = append(s.After, func(ctx context.Context) error {
			if err := _defaultDB.Close(); err != nil {
				return err
			}
			return DropTestDatabase(ctx, _defaultDB)
		})
	}
}
