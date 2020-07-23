package stats

import "github.com/blend/go-sdk/env"

// AddDefaultTagsFromEnv adds default tags to a collector from environment values.
func AddDefaultTagsFromEnv(collector Collector) {
	if collector == nil {
		return
	}
	collector.AddDefaultTags(
		Tag(TagService, env.Env().String("SERVICE_NAME")),
		Tag(TagEnv, env.Env().String("SERVICE_ENV")),
		Tag(TagContainer, env.Env().String("HOSTNAME")),
	)
}

// AddDefaultTags adds default tags to a stats collector.
func AddDefaultTags(collector Collector, serviceName, serviceEnv, container string) {
	if collector == nil {
		return
	}
	collector.AddDefaultTags(
		Tag(TagService, serviceName),
		Tag(TagEnv, serviceEnv),
		Tag(TagContainer, container),
	)
}
