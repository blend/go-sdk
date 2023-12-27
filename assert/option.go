/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
