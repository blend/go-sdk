/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
