/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
