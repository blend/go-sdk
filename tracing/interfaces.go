/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
