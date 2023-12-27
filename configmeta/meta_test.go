/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configmeta

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func Test_Meta_Resolve_envOverrides(t *testing.T) {
	its := assert.New(t)

	bareCfg := &Meta{
		ServiceName:	"not-mock-test",
		ServiceEnv:	"not-mock-test-env",
		Hostname:	"not-mock-hostname",
	}

	vars := env.Vars{
		env.VarServiceName:	"mock-test",
		env.VarServiceEnv:	"mock-test-env",
		env.VarHostname:	"mock-test-hostname",
	}
	ctx := env.WithVars(context.Background(), vars)
	err := bareCfg.Resolve(ctx)
	its.Nil(err)

	its.Equal("mock-test", bareCfg.ServiceName)
	its.Equal("mock-test-env", bareCfg.ServiceEnv)
	its.Equal("mock-test-hostname", bareCfg.Hostname)

	its.Equal("mock-test", bareCfg.ServiceNameOrDefault())
	its.Equal("mock-test-env", bareCfg.ServiceEnvOrDefault())
}
