package grpcmetrics

import (
	"context"
	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
)

// AddListeners adds grpc listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(grpcutil.RPC, stats.ListenerNameStats, grpcutil.NewRPCEventListener(func(_ context.Context, re grpcutil.RPCEvent) {
		var method string
		if len(re.Method) > 0 {
			method = stats.Tag(TagRPCMethod, re.Method)
		} else {
			method = stats.Tag(TagRPCMethod, RPCMethodUnknown)
		}

		engine := stats.Tag(TagRPCEngine, re.Engine)
		peer := stats.Tag(TagRPCPeer, re.Peer)
		tags := []string{
			method, engine, peer,
		}

		if re.Err != nil {
			tags = append(tags, stats.TagError)
		}
		collector.Increment(MetricNameRPC, tags...)
		collector.TimeInMilliseconds(MetricNameRPCElapsed, re.Elapsed, tags...)
	}))
}
