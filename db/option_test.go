package db

import (
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
	assert.Nil(OptStatementInterceptor(func(label, statement string) string { return "ok!" })(c))
	assert.NotNil(c.StatementInterceptor)

	assert.Empty(c.Config.DSN)
	assert.Nil(OptConfig(Config{DSN: "foo"})(c))
	assert.Equal("foo", c.Config.DSN)
}
