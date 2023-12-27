/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcstats

import "github.com/blend/go-sdk/grpcutil"

// Tag constants
const (
	TagRPCMethod			string	= "rpc_method"
	TagRPCPeer			string	= "rpc_peer"
	TagRPCStreamMessageDirection	string	= "rpc_stream_msg_direction"
	TagRPCEngine			string	= "rpc_peer"
	TagRPCAuthority			string	= "rpc_authority"

	RPCMethodUnknown	string	= "unknown"

	MetricNameRPC				string	= string(grpcutil.FlagRPC)
	MetricNameRPCStreamMessage		string	= string(grpcutil.FlagRPCStreamMessage)
	MetricNameRPCElapsed			string	= MetricNameRPC + ".elapsed"
	MetricNameRPCElapsedLast		string	= MetricNameRPCElapsed + ".last"
	MetricNameRPCStreamMessageElapsed	string	= MetricNameRPCStreamMessage + ".elapsed"
	MetricNameRPCStreamMessageElapsedLast	string	= MetricNameRPCStreamMessageElapsed + ".last"
)
