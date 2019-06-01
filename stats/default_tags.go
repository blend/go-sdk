package stats

import "github.com/blend/go-sdk/env"

// AddDefaultTagsFromEnv adds default tags to a collector from environment values.
func AddDefaultTagsFromEnv(collector Collector) {
	if collector == nil {
		return
	}
	collector.AddDefaultTag(TagService, env.Env().String("SERVICE_NAME"))
	collector.AddDefaultTag(TagEnv, env.Env().String("SERVICE_ENV"))
	collector.AddDefaultTag(TagContainer, env.Env().String("HOSTNAME"))
}

// AddDefaultTags adds default tags to a stats collector.
func AddDefaultTags(collector Collector, serviceName, serviceEnv, container string) {
	if collector == nil {
		return
	}
	collector.AddDefaultTag(TagService, serviceName)
	collector.AddDefaultTag(TagEnv, serviceEnv)
	collector.AddDefaultTag(TagContainer, container)
}
