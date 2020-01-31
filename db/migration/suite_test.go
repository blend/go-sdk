package migration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
)

func TestSuite_Apply(t *testing.T) {
	a := assert.New(t)
	testSchemaName := buildTestSchemaName()
	err := db.IgnoreExecResult(defaultDB().Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName)))
	a.Nil(err)
	s := New(OptLog(logger.None()), OptGroups(createTestMigrations(testSchemaName)...))
	defer func() {
		// pq can't parameterize Drop
		err := db.IgnoreExecResult(defaultDB().Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName)))
		a.Nil(err)
	}()
	err = s.Apply(context.Background(), defaultDB())
	a.Nil(err)

	ok, err := defaultDB().Query("SELECT 1 FROM pg_catalog.pg_indexes where indexname = $1 and tablename = $2", "idx_created_foo", "table_test_foo").Any()
	a.Nil(err)
	a.True(ok)

	ap, sk, fl, tot := s.Results()
	a.Equal(4, ap)
	a.Equal(1, sk)
	a.Equal(0, fl)
	a.Equal(5, tot)
}

func TestSuite_ApplyFails(t *testing.T) {
	a := assert.New(t)
	testSchemaName := buildTestSchemaName()
	err := db.IgnoreExecResult(defaultDB().Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName)))
	a.Nil(err)
	s := New(OptLog(logger.None()), OptGroups(createTestMigrations(testSchemaName)...))
	s.Groups = append(s.Groups, NewGroupWithActions(NewStep(Always(), Actions(Statements(`INSERT INTO tab_not_exists VALUES (1, 'blah', CURRENT_TIMESTAMP');`)))))
	defer func() {
		// pq can't parameterize Drop
		err := db.IgnoreExecResult(defaultDB().Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName)))
		a.Nil(err)
	}()
	err = s.Apply(context.Background(), defaultDB())
	a.NotNil(err)

	ok, err := defaultDB().Query("SELECT 1 FROM pg_catalog.pg_indexes where indexname = $1 and tablename = $2", "idx_created_foo", "table_test_foo").Any()
	a.Nil(err)
	a.True(ok)

	ap, sk, fl, tot := s.Results()
	a.Equal(4, ap)
	a.Equal(1, sk)
	a.Equal(1, fl)
	a.Equal(6, tot)
}

func createTestMigrations(testSchemaName string) []*Group {
	return []*Group{
		NewGroupWithActions(
			NewStep(
				SchemaNotExists(testSchemaName),
				Actions(
					// pq can't parameterize Create
					func(i context.Context, connection *db.Connection, tx *sql.Tx) error {
						err := db.IgnoreExecResult(connection.Exec(fmt.Sprintf("CREATE SCHEMA %s;", testSchemaName)))
						if err != nil {
							return err
						}
						return nil
					},
					// Test NoOp
					NoOp,
					func(i context.Context, connection *db.Connection, tx *sql.Tx) error {
						// This is a hack to set the schema on the connection
						(&connection.Config).Schema = testSchemaName
						return nil
					},
				))),
		NewGroupWithActions(
			NewStep(
				TableNotExists("table_test_foo"),
				Exec(fmt.Sprintf("CREATE TABLE %s.table_test_foo (id serial not null primary key, something varchar(32) not null);", testSchemaName)),
			),
			NewStep(
				ColumnNotExists("table_test_foo", "created_foo"),
				Statements(fmt.Sprintf("ALTER TABLE %s.table_test_foo ADD COLUMN created_foo timestamp not null;", testSchemaName)),
			)),
		NewGroup(OptGroupSkipTransaction(), OptGroupActions(
			NewStep(
				IndexNotExists("table_test_foo", "idx_created_foo"),
				Statements(fmt.Sprintf("CREATE INDEX CONCURRENTLY idx_created_foo ON %s.table_test_foo(created_foo);", testSchemaName)),
			))),
		NewGroupWithActions(
			NewStep(
				TableNotExists("table_test_foo"),
				Exec(fmt.Sprintf("CREATE TABLE %s.table_test_foo (id serial not null primary key, something varchar(32) not null);", testSchemaName)),
			)),
	}
}
