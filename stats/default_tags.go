/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
		Tag(TagClusterName, env.Env().String("CLUSTER_NAME")),
	)
}

// AddDefaultTags adds default tags to a stats collector.
func AddDefaultTags(collector Collector, serviceName, serviceEnv, clusterName string) {
	if collector == nil {
		return
	}
	collector.AddDefaultTags(
		Tag(TagService, serviceName),
		Tag(TagEnv, serviceEnv),
		Tag(TagClusterName, clusterName),
	)
}
