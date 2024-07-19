/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

var (
	_ configutil.Resolver = (*ViewCacheConfig)(nil)
)

func TestViewCacheConfigResolve(t *testing.T) {
	assert := assert.New(t)
	vcc := &ViewCacheConfig{}

	defer env.Restore()
	env.SetEnv(env.New())
	assert.Nil(vcc.Resolve(env.WithVars(context.Background(), env.Env())))
	assert.False(vcc.LiveReload)

	env.Env().Set("LIVE_RELOAD", "true")
	assert.Nil(vcc.Resolve(env.WithVars(context.Background(), env.Env())))
	assert.True(vcc.LiveReload)
}

func TestViewCacheConfigBufferPool(t *testing.T) {
	assert := assert.New(t)
	vcc := &ViewCacheConfig{}
	assert.Equal(DefaultViewBufferPoolSize, vcc.BufferPoolSizeOrDefault())
	vcc.BufferPoolSize = 10
	assert.Equal(vcc.BufferPoolSize, vcc.BufferPoolSizeOrDefault())
}

func TestViewCacheConfigTemplateNames(t *testing.T) {
	assert := assert.New(t)
	vcc := &ViewCacheConfig{}
	assert.Equal(DefaultTemplateNameInternalError, vcc.InternalErrorTemplateNameOrDefault())
	vcc.InternalErrorTemplateName = "hello"
	assert.Equal(vcc.InternalErrorTemplateName, vcc.InternalErrorTemplateNameOrDefault())

	assert.Equal(DefaultTemplateNameBadRequest, vcc.BadRequestTemplateNameOrDefault())
	vcc.BadRequestTemplateName = "hello"
	assert.Equal(vcc.BadRequestTemplateName, vcc.BadRequestTemplateNameOrDefault())

	assert.Equal(DefaultTemplateNameNotFound, vcc.NotFoundTemplateNameOrDefault())
	vcc.NotFoundTemplateName = "hello"
	assert.Equal(vcc.NotFoundTemplateName, vcc.NotFoundTemplateNameOrDefault())

	assert.Equal(DefaultTemplateNameNotAuthorized, vcc.NotAuthorizedTemplateNameOrDefault())
	vcc.NotAuthorizedTemplateName = "hello"
	assert.Equal(vcc.NotAuthorizedTemplateName, vcc.NotAuthorizedTemplateNameOrDefault())

	assert.Equal(DefaultTemplateNameStatus, vcc.StatusTemplateNameOrDefault())
	vcc.StatusTemplateName = "hello"
	assert.Equal(vcc.StatusTemplateName, vcc.StatusTemplateNameOrDefault())
}
