/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "context"

// Checker is a type that can be checked for SLA status.
type Checker interface {
	Check(context.Context) error
}

var (
	_ Checker = (*CheckerFunc)(nil)
)

// CheckerFunc implements Checker.
type CheckerFunc func(context.Context) error

// Check implements Checker.
func (cf CheckerFunc) Check(ctx context.Context) error {
	return cf(ctx)
}
