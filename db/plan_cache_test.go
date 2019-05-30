package db

import (
	"context"
	"testing"

	assert "github.com/blend/go-sdk/assert"
)

func TestStatementCachePrepare(t *testing.T) {
	assert := assert.New(t)

	sc := NewPlanCache().WithConnection(defaultDB().Connection)

	query := "select 'ok'"
	stmt, err := sc.PrepareContext(context.Background(), query, query)
	assert.Nil(err)
	assert.NotNil(stmt)
	assert.True(sc.HasStatement(query))

	// shoul result in cache hit
	stmt, err = sc.PrepareContext(context.Background(), query, query)
	assert.NotNil(stmt)
	assert.True(sc.HasStatement(query))
}
