/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redisstats

import "github.com/blend/go-sdk/redis"

// Metric and tag names etc.
const (
	MetricName		string	= string(redis.Flag)
	MetricNameElapsed	string	= MetricName + ".elapsed"
	MetricNameElapsedLast	string	= MetricNameElapsed + ".last"

	TagNetwork	string	= "network"
	TagAddr		string	= "addr"
	TagDB		string	= "db"
	TagOp		string	= "op"
)
