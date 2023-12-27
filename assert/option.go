/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package assert

import "io"

// Option mutates assertions.
type Option func(*Assertions)

// OptOutput sets the output for assertions.
func OptOutput(wr io.Writer) Option {
	return func(a *Assertions) {
		a.Output = wr
	}
}
