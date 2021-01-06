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
