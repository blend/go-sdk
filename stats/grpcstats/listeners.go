package grpcstats

import (
	"context"
	"strconv"

	"google.golang.org/grpc/status"

	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

func getErrorTag(err error) string {
	if e, ok := status.FromError(err); ok {
		code := e.Code()
		return stats.Tag(stats.TagError, strconv.Itoa(int(code)))
	}
	return stats.TagError
}

// AddListeners adds grpc listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(grpcutil.FlagRPC, stats.ListenerNameStats, grpcutil.NewRPCEventListener(func(ctx context.Context, re grpcutil.RPCEvent) {
		var tags []string

		labels := logger.GetLabels(ctx)
		for key, value := range labels {
			tags = append(tags, stats.Tag(key, value))
		}

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
			tags = append(tags, getErrorTag(re.Err))
		}
		_ = collector.Increment(MetricNameRPC, tags...)
		_ = collector.Gauge(MetricNameRPCElapsed, timeutil.Milliseconds(re.Elapsed), tags...)
		_ = collector.TimeInMilliseconds(MetricNameRPCElapsed, re.Elapsed, tags...)
		_ = collector.Distribution(MetricNameRPCElapsed, timeutil.Milliseconds(re.Elapsed), tags...)
	}))

	log.Listen(grpcutil.FlagRPCStreamMessage, stats.ListenerNameStats, grpcutil.NewRPCStreamMessageEventListener(func(ctx context.Context, re grpcutil.RPCStreamMessageEvent) {
		var tags []string

		labels := logger.GetLabels(ctx)
		for key, value := range labels {
			tags = append(tags, stats.Tag(key, value))
		}

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

		if re.Direction != "" {
			tags = append(tags, stats.Tag(TagRPCStreamMessageDirection, string(re.Direction)))
		}
		if re.Err != nil {
			tags = append(tags, getErrorTag(re.Err))
		}
		_ = collector.Increment(MetricNameRPCStreamMessage, tags...)
		_ = collector.Gauge(MetricNameRPCStreamMessageElapsed, timeutil.Milliseconds(re.Elapsed), tags...)
		_ = collector.TimeInMilliseconds(MetricNameRPCStreamMessageElapsed, re.Elapsed, tags...)
		_ = collector.Distribution(MetricNameRPCStreamMessageElapsed, timeutil.Milliseconds(re.Elapsed), tags...)
	}))
}
