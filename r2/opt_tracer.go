/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

// OptTracer sets the optional trace handler.
func OptTracer(tracer Tracer) Option {
	return func(r *Request) error {
		r.Tracer = tracer
		return nil
	}
}
