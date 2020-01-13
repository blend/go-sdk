package db

import (
	"context"
	"database/sql"
	"sync"
	"github.com/blend/go-sdk/ex"
)

// NewPlanCache returns a new `PlanCache`.
func NewPlanCache(conn *sql.DB) *PlanCache {
	return &PlanCache{
		Connection: conn,
		Cache:      sync.Map{},
	}
}

// PlanCache is a cache of prepared statements.
type PlanCache struct {
	Connection *sql.DB
	Cache      sync.Map
}

// Close implements io.Closer.
// It ranges over the cached statements and closes them.
func (pc *PlanCache) Close() (err error) {
	pc.Cache.Range(func(k, v interface{}) bool {
		err = v.(*sql.Stmt).Close()
		return err == nil
	})
	return
}

// Has returns if the cache contains a statement.
func (pc *PlanCache) Has(key string) bool {
	_, hasStmt := pc.Cache.Load(key)
	return hasStmt
}

// Invalidate removes a statement from the cache.
func (pc *PlanCache) Invalidate(key string) (err error) {
	stmt, ok := pc.Cache.Load(key)
	if !ok {
		return
	}
	pc.Cache.Delete(key)
	return stmt.(*sql.Stmt).Close()
}

// PrepareContext returns a cached expression for a statement, or creates and caches a new one.
func (pc *PlanCache) PrepareContext(ctx context.Context, key, statement string) (*sql.Stmt, error) {
	if key == "" {
		return nil, ex.New(ErrPlanCacheKeyUnset)
	}

	if stmt, hasStmt := pc.Cache.Load(key); hasStmt {
		return stmt.(*sql.Stmt), nil
	}

	stmt, err := pc.Connection.PrepareContext(ctx, statement)
	if err != nil {
		return nil, ex.New(err)
	}

	pc.Cache.Store(key, stmt)
	return stmt, nil
}
