/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package traceserver

// Headers
const (
	// HeaderTraceCount is a header containing the number of traces in the payload
	HeaderTraceCount	= "X-Datadog-Trace-Count"
	HeaderContainerID	= "Datadog-Container-ID"
)

// ContentTypes
const (
	ContentTypeApplicationMessagePack = "application/msgpack"
)
