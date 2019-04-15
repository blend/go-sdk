package db

import (
	"database/sql"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

// TestConnectionSanityCheck tests if we can connect to the db, a.k.a., if the underlying driver works.
func TestConnectionUseBeforeOpen(t *testing.T) {
	assert := assert.New(t)

	conn, err := New()
	assert.Nil(err)

	tx, err := conn.Begin()
	assert.NotNil(err)
	assert.True(ex.Is(ErrConnectionClosed, err))
	assert.Nil(tx)
}

// TestConnectionSanityCheck tests if we can connect to the db, a.k.a., if the underlying driver works.
func TestConnectionSanityCheck(t *testing.T) {
	assert := assert.New(t)

	conn, err := New(OptConfigFromEnv())
	assert.Nil(err)
	str := conn.Config.CreateDSN()
	_, err = sql.Open("postgres", str)
	assert.Nil(err)
}

func TestPrepare(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	a.Nil(err)
}

func TestQuery(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = seedObjects(100, tx)
	a.Nil(err)

	objs := []benchObj{}
	err = Default().Invoke(OptTx(tx)).Query("select * from bench_object").OutMany(&objs)
	a.Nil(err)
	a.NotEmpty(objs)
}

func TestConnectionStatementCacheExecute(t *testing.T) {
	a := assert.New(t)

	conn, err := New(OptConfigFromEnv())
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()
	conn.PlanCache.WithEnabled(true)

	a.Nil(conn.Exec("select 'ok!'"))
	a.Nil(conn.Exec("select 'ok!'"))
	a.False(conn.PlanCache.HasStatement("select 'ok!'"))

	a.Nil(conn.Invoke(OptCachedPlanKey("ping")).Exec("select 'ok!'"))
	a.Nil(conn.Invoke(OptCachedPlanKey("ping")).Exec("select 'ok!'"))
	a.True(conn.PlanCache.HasStatement("ping"))
}

func TestConnectionStatementCacheQuery(t *testing.T) {
	a := assert.New(t)

	conn, err := New(OptConfigFromEnv())
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()

	conn.PlanCache.WithEnabled(true)

	var ok string
	a.Nil(conn.Invoke(OptCachedPlanKey("status")).Query("select 'ok!'").Scan(&ok))
	a.Equal("ok!", ok)

	a.Nil(conn.Invoke(OptCachedPlanKey("status")).Query("select 'ok!'").Scan(&ok))
	a.Equal("ok!", ok)

	a.True(conn.PlanCache.HasStatement("status"))
}

func TestConnectionOpen(t *testing.T) {
	a := assert.New(t)

	conn, err := New(OptConfigFromEnv())
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()

	a.NotNil(conn.BufferPool)
	a.NotNil(conn.Connection)
	a.NotNil(conn.PlanCache)
}

func TestExec(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = Default().Invoke(OptTx(tx)).Exec("select 'ok!'")
	a.Nil(err)
}

func TestConnectionInvalidatesBadCachedStatements(t *testing.T) {
	assert := assert.New(t)

	conn, err := New(OptConfigFromEnv())
	assert.Nil(err)
	assert.Nil(conn.Open())
	defer conn.Close()

	conn.PlanCache.WithEnabled(true)

	createTableStatement := `CREATE TABLE state_invalidation (id int not null, name varchar(64))`
	insertStatement := `INSERT INTO state_invalidation (id, name) VALUES ($1, $2)`
	alterTableStatement := `ALTER TABLE state_invalidation ALTER COLUMN id TYPE bigint;`
	dropTableStatement := `DROP TABLE state_invalidation`
	queryStatement := `SELECT * from state_invalidation`

	defer func() {
		err = conn.Exec(dropTableStatement)
		assert.Nil(err)
	}()

	err = conn.Exec(createTableStatement)
	assert.Nil(err)

	err = conn.Exec(insertStatement, 1, "Foo")
	assert.Nil(err)

	err = conn.Exec(insertStatement, 2, "Bar")
	assert.Nil(err)

	_, err = conn.Query(queryStatement).Any()
	assert.Nil(err)

	err = conn.Exec(alterTableStatement)
	assert.Nil(err)

	// normally this would result in a busted cached query plan.
	// we need to invalidate the cache and make this work.
	_, err = conn.Query(queryStatement).Any()
	assert.Nil(err)

	_, err = conn.Query(queryStatement).Any()
	assert.Nil(err)
}

// TestConnectionConfigSetsDatabase tests if we set the .database property on open.
func TestConnectionConfigSetsDatabase(t *testing.T) {
	assert := assert.New(t)
	conn, err := New(OptConfigFromEnv())
	assert.Nil(err)
	assert.Nil(conn.Open())
	defer conn.Close()
	assert.NotEmpty(conn.Config.DatabaseOrDefault())
}
