/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"context"

	"github.com/blend/go-sdk/sanitize"
)

// NewAddListenerOptions creates a new add listener options.
func NewAddListenerOptions(opts ...AddListenerOption) AddListenerOptions {
	var options AddListenerOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// AddListenerOptions are options for adding listeners.
type AddListenerOptions struct {
	IncludeLoggerLabelsAsTags	bool
	RequestSanitizeDefaults		[]sanitize.RequestOption
	DefaultTags			map[string]string
}

// GetLoggerLabelsAsTags gets the logger tags from a context if they're set to be included.
func (options AddListenerOptions) GetLoggerLabelsAsTags(ctx context.Context) (tags []string) {
	if options.IncludeLoggerLabelsAsTags {
		tags = GetLoggerLabelsAsTags(ctx)
	}
	return
}

// OptIncludeLoggerLabelsAsTags includes logger labels as tags.
func OptIncludeLoggerLabelsAsTags(include bool) AddListenerOption {
	return func(a *AddListenerOptions) { a.IncludeLoggerLabelsAsTags = include }
}

// OptRequestSanitizeDefaults includes logger labels as tags.
func OptRequestSanitizeDefaults(opts ...sanitize.RequestOption) AddListenerOption {
	return func(a *AddListenerOptions) { a.RequestSanitizeDefaults = opts }
}

// GetDefaultsAsTags gets a set of constant tags at collector initialization to be included.
func (options AddListenerOptions) GetDefaultsAsTags() (tags []string) {
	if len(options.DefaultTags) > 0 {
		tags = GetDefaultsAsTags(options.DefaultTags)
	}
	return
}

// OptRequestIncludeConstantsAsTags includes fixed tags.
func OptRequestIncludeConstantsAsTags(opts map[string]string) AddListenerOption {
	return func(a *AddListenerOptions) { a.DefaultTags = opts }
}

// AddListenerOption mutates AddListenerOptions
type AddListenerOption func(*AddListenerOptions)
