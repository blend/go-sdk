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
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestInvocationOptions(t *testing.T) {
	assert := assert.New(t)

	i := &Invocation{}

	assert.Empty(i.Label)
	OptLabel("label")(i)
	assert.Equal("label", i.Label)

	assert.Nil(i.StatementInterceptor)
	OptInvocationStatementInterceptor(func(_ context.Context, label, statement string) (string, error) { return "OK!", nil })(i)
	assert.NotNil(i.StatementInterceptor)

	assert.Nil(i.Context)
	OptContext(context.Background())(i)
	assert.NotNil(i.Context)

	assert.Nil(i.Cancel)
	OptCancel(func() {})(i)
	assert.NotNil(i.Cancel)

	i.Cancel = nil
	assert.Nil(i.Cancel)
	OptTimeout(5 * time.Second)(i)
	assert.NotNil(i.Cancel)
	assert.NotNil(i.Context)

	i.DB = defaultDB().Connection
	assert.NotNil(i.DB)
	OptTx(nil)(i)
	assert.NotNil(i.DB)

	i.DB = nil
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	OptTx(tx)(i)
	assert.NotNil(i.DB)

	i.DB = nil
	OptInvocationDB(defaultDB().Connection)(i)
	assert.NotNil(i.DB)

	i.DB = nil
	OptInvocationDB(tx)(i)
	assert.NotNil(i.DB)

	i.StatementInterceptor = nil
	OptInvocationStatementInterceptor(failInterceptor)(i)
	assert.NotNil(i.StatementInterceptor)
}
