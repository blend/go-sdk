package db

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPlanCache(t *testing.T) {
	assert := assert.New(t)

	sc := NewPlanCache(defaultDB().Connection)

	query := "select 'ok'"
	stmt, err := sc.PrepareContext(context.Background(), query, query)
	assert.Nil(err)
	assert.NotNil(stmt)
	assert.True(sc.Has(query))

	// shoul result in cache hit
	stmt, err = sc.PrepareContext(context.Background(), query, query)
	assert.NotNil(stmt)
	assert.True(sc.Has(query))

	sc.Invalidate(query)
	assert.False(sc.Has(query))
	assert.Nil(sc.Close())
}
