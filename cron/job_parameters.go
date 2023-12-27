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
)

// JobParameters is a loose association to map[string]string.
type JobParameters = map[string]string

type contextKeyJobParameters struct{}

// WithJobParameterValues adds job invocation parameter values to a context.
func WithJobParameterValues(ctx context.Context, values JobParameters) context.Context {
	return context.WithValue(ctx, contextKeyJobParameters{}, values)
}

// GetJobParameterValues gets parameter values from a given context.
func GetJobParameterValues(ctx context.Context) JobParameters {
	if value := ctx.Value(contextKeyJobParameters{}); value != nil {
		if typed, ok := value.(JobParameters); ok {
			return typed
		}
	}
	return nil
}

// MergeJobParameterValues merges values from many sources.
// The order is important for which value set's keys take precedence.
func MergeJobParameterValues(values ...JobParameters) JobParameters {
	output := make(JobParameters)
	for _, set := range values {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}
