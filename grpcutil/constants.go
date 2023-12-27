/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

// MetaTags
// These are common tags found in the metadata for rpc calls, both unary and streaming.
const (
	MetaTagAuthority   = "authority"
	MetaTagContentType = "content-type"
	MetaTagUserAgent   = "user-agent"
)

// Our default engine
const (
	EngineGRPC = "grpc"
)
