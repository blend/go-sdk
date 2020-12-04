package grpcstats

import (
	"context"

	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds grpc listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(grpcutil.FlagRPC, stats.ListenerNameStats, grpcutil.NewRPCEventListener(func(_ context.Context, re grpcutil.RPCEvent) {
		var tags []string

		if len(re.Method) > 0 {
			tags = append(tags, stats.Tag(TagRPCMethod, re.Method))
		} else {
			tags = append(tags, stats.Tag(TagRPCMethod, RPCMethodUnknown))
		}

		if re.Engine != "" {
			tags = append(tags, stats.Tag(TagRPCEngine, re.Engine))
		}
		if re.Peer != "" {
			tags = append(tags, stats.Tag(TagRPCPeer, re.Peer))
		}

		if re.Err != nil {
			tags = append(tags, stats.TagError)
		}
		_ = collector.Increment(MetricNameRPC, tags...)
		_ = collector.Gauge(MetricNameRPCElapsed, timeutil.Milliseconds(re.Elapsed), tags...)
		_ = collector.TimeInMilliseconds(MetricNameRPCElapsed, re.Elapsed, tags...)
	}))
}
