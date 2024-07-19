/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
