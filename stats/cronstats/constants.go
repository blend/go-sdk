/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cronstats

// HTTP stats constants
const (
	MetricNameCron            = "cron.job"
	MetricNameCronElapsed     = MetricNameCron + ".elapsed"
	MetricNameCronElapsedLast = MetricNameCronElapsed + ".last"

	TagJob       = "job"
	TagJobStatus = "job_status"
)
