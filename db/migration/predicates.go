package migration

import (
	"database/sql"
	"strings"
	"github.com/blend/go-sdk/db"
)

// PredicateTableExists returns if a table exists in the default schema of the given connection.
func PredicateTableExists(c *db.Connection, tx *sql.Tx, tableName string) (bool, error) {
	return PredicateTableExistsInSchema(c, tx, c.Config.SchemaOrDefault(), tableName)
}

// PredicateTableExistsInSchema returns if a table exists in a specific schema on the given connection.
func PredicateTableExistsInSchema(c *db.Connection, tx *sql.Tx, schemaName, tableName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM pg_catalog.pg_tables WHERE tablename = $1 AND schemaname = $2`,
		strings.ToLower(tableName), strings.ToLower(schemaName)).Any()
}

// PredicateColumnExists returns if a column exists on a table in the default schema of the given connection.
func PredicateColumnExists(c *db.Connection, tx *sql.Tx, tableName, columnName string) (bool, error) {
	return PredicateColumnExistsInSchema(c, tx, c.Config.SchemaOrDefault(), tableName, columnName)
}

// PredicateColumnExistsInSchema returns if a column exists on a table in a specific schema on the given connection.
func PredicateColumnExistsInSchema(c *db.Connection, tx *sql.Tx, schemaName, tableName, columnName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(
		`SELECT 1 FROM information_schema.columns WHERE column_name = $1 AND table_name = $2 AND table_schema = $3`,
		strings.ToLower(columnName), strings.ToLower(tableName), strings.ToLower(schemaName)).Any()
}

// PredicateConstraintExists returns if a constraint exists on a table in the default schema of the given connection.
func PredicateConstraintExists(c *db.Connection, tx *sql.Tx, tableName, constraintName string) (bool, error) {
	return PredicateConstraintExistsInSchema(c, tx, c.Config.SchemaOrDefault(), tableName, constraintName)
}

// PredicateConstraintExistsInSchema returns if a constraint exists on a table in a specific schema on the given connection.
func PredicateConstraintExistsInSchema(c *db.Connection, tx *sql.Tx, schemaName, tableName, constraintName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(
		`SELECT 1 FROM information_schema.constraint_column_usage WHERE constraint_name = $1 AND table_name = $2 AND table_schema = $3`,
		strings.ToLower(constraintName), strings.ToLower(tableName), strings.ToLower(schemaName)).Any()
}

// PredicateIndexExists returns if a index exists on a table in the default schema of the given connection.
func PredicateIndexExists(c *db.Connection, tx *sql.Tx, tableName, indexName string) (bool, error) {
	return PredicateIndexExistsInSchema(c, tx, c.Config.SchemaOrDefault(), tableName, indexName)
}

// PredicateIndexExistsInSchema returns if a index exists on a table in a specific schema on the given connection.
func PredicateIndexExistsInSchema(c *db.Connection, tx *sql.Tx, schemaName, tableName, indexName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(
		`SELECT 1 FROM pg_catalog.pg_indexes where indexname = $1 and tablename = $2 AND schemaname = $3`,
		strings.ToLower(indexName), strings.ToLower(tableName), strings.ToLower(schemaName)).Any()
}

// PredicateRoleExists returns if a role exists or not.
func PredicateRoleExists(c *db.Connection, tx *sql.Tx, roleName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM pg_catalog.pg_roles WHERE rolname ilike $1`, roleName).Any()
}
