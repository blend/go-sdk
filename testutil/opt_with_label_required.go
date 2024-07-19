/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"context"

	"github.com/blend/go-sdk/db"
)

// OptWithStatementLabelRequired adds a defaultdb interceptor that enforces
// that statement labels must be present on all statements.
func OptWithStatementLabelRequired() Option {
	return func(s *Suite) {
		s.Before = append(s.Before, func(ctx context.Context) error {
			_defaultDB.StatementInterceptor = db.StatementInterceptorChain(
				_defaultDB.StatementInterceptor,
				db.LabelRequiredStatementInterceptor,
			)
			return nil
		})
	}
}
