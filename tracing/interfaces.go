/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package tracing

// SpanIDProvider is a tracing span context that has a SpanID getter
type SpanIDProvider interface {
	SpanID() uint64
}

// TraceIDProvider is a tracing span context that has a TraceID getter
type TraceIDProvider interface {
	TraceID() uint64
}
