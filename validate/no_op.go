/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

// NoOp is an empty validated.
type NoOp struct{}

// Validate implements the no op.
func (no NoOp) Validate() error	{ return nil }
