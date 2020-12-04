package grpcstats

import "github.com/blend/go-sdk/grpcutil"

// Tag constants
const (
	TagRPCMethod    string = "rpc_method"
	TagRPCPeer      string = "rpc_peer"
	TagRPCEngine    string = "rpc_peer"
	TagRPCAuthority string = "rpc_authority"

	RPCMethodUnknown string = "unknown"

	MetricNameRPC        string = string(grpcutil.FlagRPC)
	MetricNameRPCElapsed string = MetricNameRPC + ".elapsed"
)
