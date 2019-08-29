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
	OptInvocationStatementInterceptor(func(label, statement string) string { return "OK!" })(i)
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
	OptDB(defaultDB().Connection)(i)
	assert.NotNil(i.DB)

	i.DB = nil
	OptDB(tx)(i)
	assert.NotNil(i.DB)
}
