/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

// OptTracer sets the optional trace handler.
func OptTracer(tracer Tracer) Option {
	return func(r *Request) error {
		r.Tracer = tracer
		return nil
	}
}
