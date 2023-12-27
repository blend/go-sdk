/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"context"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/ref"
)

var (
	_ configutil.Resolver = (*JobConfig)(nil)
)

// JobConfig is a configuration set for a job.
type JobConfig struct {
	// Disabled determines if the job should be automatically scheduled or not.
	Disabled	*bool	`json:"disabled" yaml:"disabled"`
	// Description is an optional string to describe what the job does.
	Description	string	`json:"description" yaml:"description"`
	// Labels define extra metadata that can be used to filter jobs.
	Labels	map[string]string	`json:"labels" yaml:"labels"`
	// ParameterValues act as default parameters for a given job.
	ParameterValues	JobParameters	`json:"parameterValues" yaml:"parameterValues"`
	// Timeout represents the abort threshold for the job.
	Timeout	time.Duration	`json:"timeout" yaml:"timeout"`
	// ShutdownGracePeriod represents the time a job is given to clean itself up.
	ShutdownGracePeriod	time.Duration	`json:"shutdownGracePeriod" yaml:"shutdownGracePeriod"`
	// SkipLoggerTrigger skips triggering logger events if it is set to true.
	SkipLoggerTrigger	bool	`json:"skipLoggerTrigger" yaml:"skipLoggerTrigger"`
}

// Resolve implements configutil.Resolver.
func (jc *JobConfig) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetBoolPtr(&jc.Disabled, configutil.Bool(jc.Disabled), configutil.Bool(ref.Bool(DefaultDisabled))),
		configutil.SetDuration(&jc.Timeout, configutil.Duration(jc.Timeout), configutil.Duration(DefaultTimeout)),
		configutil.SetDuration(&jc.ShutdownGracePeriod, configutil.Duration(jc.ShutdownGracePeriod), configutil.Duration(DefaultShutdownGracePeriod)),
	)
}

// DisabledOrDefault returns a value or a default.
func (jc JobConfig) DisabledOrDefault() bool {
	if jc.Disabled != nil {
		return *jc.Disabled
	}
	return DefaultDisabled
}

// TimeoutOrDefault returns a value or a default.
func (jc JobConfig) TimeoutOrDefault() time.Duration {
	if jc.Timeout > 0 {
		return jc.Timeout
	}
	return DefaultTimeout
}

// ShutdownGracePeriodOrDefault returns a value or a default.
func (jc JobConfig) ShutdownGracePeriodOrDefault() time.Duration {
	if jc.ShutdownGracePeriod > 0 {
		return jc.ShutdownGracePeriod
	}
	return DefaultShutdownGracePeriod
}
