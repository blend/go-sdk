/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"context"
)

// StatementInterceptor is an interceptor for statements.
type StatementInterceptor func(ctx context.Context, label, statement string) (string, error)
