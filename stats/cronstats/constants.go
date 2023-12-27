/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cronstats

// HTTP stats constants
const (
	MetricNameCron			= "cron.job"
	MetricNameCronElapsed		= MetricNameCron + ".elapsed"
	MetricNameCronElapsedLast	= MetricNameCronElapsed + ".last"

	TagJob		= "job"
	TagJobStatus	= "job_status"
)
