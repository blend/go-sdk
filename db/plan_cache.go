package db

import (
	"context"
	"database/sql"
	"sync"

	"github.com/blend/go-sdk/ex"
)

// NewPlanCache returns a new `PlanCache`.
func NewPlanCache() *PlanCache {
	return &PlanCache{
		enabled: true,
		cache:   sync.Map{},
	}
}

// PlanCache is a cache of prepared statements.
type PlanCache struct {
	conn    *sql.DB
	enabled bool
	cache   sync.Map
}

// WithConnection sets the statement cache connection.
func (pc *PlanCache) WithConnection(conn *sql.DB) *PlanCache {
	pc.conn = conn
	return pc
}

// Connection returns the underlying connection.
func (pc *PlanCache) Connection() *sql.DB {
	return pc.conn
}

// WithEnabled sets if the cache is enabled.
func (pc *PlanCache) WithEnabled(enabled bool) *PlanCache {
	pc.enabled = enabled
	return pc
}

// Enabled returns if the statement cache is enabled.
func (pc *PlanCache) Enabled() bool {
	return pc.enabled
}

// Close implements io.Closer.
func (pc *PlanCache) Close() (err error) {
	pc.cache.Range(func(k, v interface{}) bool {
		err = v.(*sql.Stmt).Close()
		return err == nil
	})
	return
}

// HasStatement returns if the cache contains a statement.
func (pc *PlanCache) HasStatement(planCacheKey string) bool {
	_, hasStmt := pc.cache.Load(planCacheKey)
	return hasStmt
}

// InvalidateStatement removes a statement from the cache.
func (pc *PlanCache) InvalidateStatement(planCacheKey string) (err error) {
	stmt, ok := pc.cache.Load(planCacheKey)
	if !ok {
		return
	}
	pc.cache.Delete(planCacheKey)
	return stmt.(*sql.Stmt).Close()
}

// PrepareContext returns a cached expression for a statement, or creates and caches a new one.
func (pc *PlanCache) PrepareContext(context context.Context, planCacheKey, statement string) (*sql.Stmt, error) {
	if len(planCacheKey) == 0 {
		return nil, ex.New(ErrPlanCacheKeyUnset)
	}

	if stmt, hasStmt := pc.cache.Load(planCacheKey); hasStmt {
		return stmt.(*sql.Stmt), nil
	}

	stmt, err := pc.conn.PrepareContext(context, statement)
	if err != nil {
		return nil, err
	}

	pc.cache.Store(planCacheKey, stmt)
	return stmt, nil
}
