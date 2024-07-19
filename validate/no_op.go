/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

// NoOp is an empty validated.
type NoOp struct{}

// Validate implements the no op.
func (no NoOp) Validate() error { return nil }
