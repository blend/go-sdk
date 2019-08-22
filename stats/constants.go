package stats

import (
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// MetricNames are names we use when sending data to the collectors.
const (
	MetricNameHTTPRequest        string = string(webutil.HTTPRequest)
	MetricNameHTTPRequestElapsed string = MetricNameHTTPRequest + ".elapsed"
	MetricNameDBQuery            string = string(db.QueryFlag)
	MetricNameDBQueryElapsed     string = MetricNameDBQuery + ".elapsed"
	MetricNameRPC                string = string(grpcutil.RPC)
	MetricNameRPCElapsed         string = MetricNameRPC + ".elapsed"
	MetricNameError              string = string(logger.Error)
)

// Tag names are names for tags, either on metrics or traces.
const (
	TagService   string = "service"
	TagJob       string = "job"
	TagEnv       string = "env"
	TagHostname  string = "hostname"
	TagContainer string = "container"

	TagRoute  string = "route"
	TagMethod string = "method"
	TagStatus string = "status"

	TagQuery    string = "query"
	TagEngine   string = "engine"
	TagDatabase string = "database"

	TagSeverity string = "severity"
	TagError    string = "error"
	TagClass    string = "class"

	TagRPCMethod    string = "rpc_method"
	TagRPCPeer      string = "rpc_peer"
	TagRPCEngine    string = "rpc_peer"
	TagRPCAuthority string = "rpc_authority"
)

// Specialized / default values
const (
	RPCMethodUnknown  string = "unknown"
	RouteNotFound     string = "not_found"
	ListenerNameStats string = "stats"
)

// Tag creates a new tag.
func Tag(key, value string) string {
	return key + ":" + value
}
