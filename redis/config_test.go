/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis_test

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/redis"
)

func Test_Config_Resolve_bare(t *testing.T) {
	its := assert.New(t)

	cfg := new(redis.Config)
	its.Nil(cfg.Resolve(context.Background()))
	its.Equal(redis.DefaultNetwork, cfg.Network)
	its.Equal(redis.DefaultAddr, cfg.Addr)
	its.Equal(redis.DefaultConnectTimeout, cfg.ConnectTimeout)
	its.Equal(redis.DefaultTimeout, cfg.Timeout)
}
