/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestOptions(t *testing.T) {
	assert := assert.New(t)

	c := &Connection{}

	assert.Nil(c.Connection)
	assert.Nil(OptConnection(&sql.DB{})(c))
	assert.NotNil(c.Connection)

	assert.Nil(c.Log)
	assert.Nil(OptLog(logger.None())(c))
	assert.NotNil(c)

	assert.Nil(c.Tracer)
	assert.Nil(OptTracer(mockTracer{})(c))
	assert.NotNil(c.Tracer)

	assert.Nil(c.StatementInterceptor)
	assert.Nil(OptStatementInterceptor(func(_ context.Context, label, statement string) (string, error) { return "ok!", nil })(c))
	assert.NotNil(c.StatementInterceptor)

	assert.Empty(c.Config.DSN)
	assert.Nil(OptConfig(Config{DSN: "foo"})(c))
	assert.Equal("foo", c.Config.DSN)
}
