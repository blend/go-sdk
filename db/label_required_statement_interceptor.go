/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"context"

	"github.com/blend/go-sdk/ex"
)

// Errors
var (
	ErrStatementLabelRequired ex.Class = "statement interceptor; a statement label is required and none was provided, cannot continue"
)

// LabelRequiredStatementInterceptor returns a statement interceptor that requires a label to be set.
func LabelRequiredStatementInterceptor(ctx context.Context, label, statement string) (string, error) {
	if label == "" {
		return statement, ErrStatementLabelRequired
	}
	return statement, nil
}
