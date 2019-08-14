package stats

import (
	"context"

	"github.com/blend/go-sdk/logger"
)

// AddRPCListeners adds rpc listeners.
func AddRPCListeners(log logger.Listenable, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.RPC, ListenerNameStats, logger.NewRPCEventListener(func(_ context.Context, re logger.RPCEvent) {
		var method string
		if len(re.Method) > 0 {
			method = Tag(TagRPCMethod, re.Method)
		} else {
			method = Tag(TagRPCMethod, RPCMethodUnknown)
		}

		engine := Tag(TagRPCEngine, re.Engine)
		peer := Tag(TagRPCPeer, re.Peer)
		tags := []string{
			method, engine, peer,
		}

		if re.Err != nil {
			tags = append(tags, TagError)
		}
		stats.Increment(MetricNameRPC, tags...)
		stats.TimeInMilliseconds(MetricNameRPCElapsed, re.Elapsed, tags...)
	}))
}
