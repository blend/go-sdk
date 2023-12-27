/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"

	"github.com/blend/go-sdk/env"
)

// ViewCacheConfig is a config for the view cache.
type ViewCacheConfig struct {
	// LiveReload indicates if we should store compiled views in memory for re-use (default), or read them from disk each load.
	LiveReload bool `json:"liveReload,omitempty" yaml:"liveReload,omitempty" env:"LIVE_RELOAD"`
	// Paths are a list of view paths to include in the templates list.
	Paths []string `json:"paths,omitempty" yaml:"paths,omitempty"`
	// BufferPoolSize is the size of the re-usable buffer pool for rendering views.
	BufferPoolSize int `json:"bufferPoolSize,omitempty" yaml:"bufferPoolSize,omitempty"`

	// InternalErrorTemplateName is the template name to use for the view result provider `InternalError` result.
	InternalErrorTemplateName string `json:"internalErrorTemplateName,omitempty" yaml:"internalErrorTemplateName,omitempty"`
	// BadRequestTemplateName is the template name to use for the view result provider `BadRequest` result.
	BadRequestTemplateName string `json:"badRequestTemplateName,omitempty" yaml:"badRequestTemplateName,omitempty"`
	// NotFoundTemplateName is the template name to use for the view result provider `NotFound` result.
	NotFoundTemplateName string `json:"notFoundTemplateName,omitempty" yaml:"notFoundTemplateName,omitempty"`
	// NotAuthorizedTemplateName is the template name to use for the view result provider `NotAuthorized` result.
	NotAuthorizedTemplateName string `json:"notAuthorizedTemplateName,omitempty" yaml:"notAuthorizedTemplateName,omitempty"`
	// StatusTemplateName is the template name to use for the view result provider status result.
	StatusTemplateName string `json:"statusTemplateName,omitempty" yaml:"statusTemplateName,omitempty"`
}

// Resolve adds extra resolution steps when we setup the config.
func (vcc *ViewCacheConfig) Resolve(ctx context.Context) error {
	return env.GetVars(ctx).ReadInto(vcc)
}

// BufferPoolSizeOrDefault gets the buffer pool size or a default.
func (vcc ViewCacheConfig) BufferPoolSizeOrDefault() int {
	if vcc.BufferPoolSize > 0 {
		return vcc.BufferPoolSize
	}
	return DefaultViewBufferPoolSize
}

// InternalErrorTemplateNameOrDefault returns the internal error template name for the app.
func (vcc ViewCacheConfig) InternalErrorTemplateNameOrDefault() string {
	if vcc.InternalErrorTemplateName != "" {
		return vcc.InternalErrorTemplateName
	}
	return DefaultTemplateNameInternalError
}

// BadRequestTemplateNameOrDefault returns the bad request template name for the app.
func (vcc ViewCacheConfig) BadRequestTemplateNameOrDefault() string {
	if vcc.BadRequestTemplateName != "" {
		return vcc.BadRequestTemplateName
	}
	return DefaultTemplateNameBadRequest
}

// NotFoundTemplateNameOrDefault returns the not found template name for the app.
func (vcc ViewCacheConfig) NotFoundTemplateNameOrDefault() string {
	if vcc.NotFoundTemplateName != "" {
		return vcc.NotFoundTemplateName
	}
	return DefaultTemplateNameNotFound
}

// NotAuthorizedTemplateNameOrDefault returns the not authorized template name for the app.
func (vcc ViewCacheConfig) NotAuthorizedTemplateNameOrDefault() string {
	if vcc.NotAuthorizedTemplateName != "" {
		return vcc.NotAuthorizedTemplateName
	}
	return DefaultTemplateNameNotAuthorized
}

// StatusTemplateNameOrDefault returns the not authorized template name for the app.
func (vcc ViewCacheConfig) StatusTemplateNameOrDefault() string {
	if vcc.StatusTemplateName != "" {
		return vcc.StatusTemplateName
	}
	return DefaultTemplateNameStatus
}
