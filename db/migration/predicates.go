/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/stringutil"
)

// PredicateTableExists returns if a table exists in the default schema of the given connection.
func PredicateTableExists(ctx context.Context, c *db.Connection, tx *sql.Tx, tableName string) (bool, error) {
	return PredicateTableExistsInSchema(ctx, c, tx, c.Config.SchemaOrDefault(), tableName)
}

// PredicateTableExistsInSchema returns if a table exists in a specific schema on the given connection.
func PredicateTableExistsInSchema(ctx context.Context, c *db.Connection, tx *sql.Tx, schemaName, tableName string) (bool, error) {
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(
		`SELECT 1 FROM pg_catalog.pg_tables WHERE tablename = $1 AND schemaname = $2`,
		tableName,
		schemaName,
	).Any()
}

// PredicateColumnExists returns if a column exists on a table in the default schema of the given connection.
func PredicateColumnExists(ctx context.Context, c *db.Connection, tx *sql.Tx, tableName, columnName string) (bool, error) {
	return PredicateColumnExistsInSchema(ctx, c, tx, c.Config.SchemaOrDefault(), tableName, columnName)
}

// PredicateColumnExistsInSchema returns if a column exists on a table in a specific schema on the given connection.
func PredicateColumnExistsInSchema(ctx context.Context, c *db.Connection, tx *sql.Tx, schemaName, tableName, columnName string) (bool, error) {
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(
		`SELECT 1 FROM information_schema.columns WHERE column_name = $1 AND table_name = $2 AND table_schema = $3`,
		columnName,
		tableName,
		schemaName,
	).Any()
}

// PredicateConstraintExists returns if a constraint exists on a table in the default schema of the given connection.
func PredicateConstraintExists(ctx context.Context, c *db.Connection, tx *sql.Tx, tableName, constraintName string) (bool, error) {
	return PredicateConstraintExistsInSchema(ctx, c, tx, c.Config.SchemaOrDefault(), tableName, constraintName)
}

// PredicateConstraintExistsInSchema returns if a constraint exists on a table in a specific schema on the given connection.
func PredicateConstraintExistsInSchema(ctx context.Context, c *db.Connection, tx *sql.Tx, schemaName, tableName, constraintName string) (bool, error) {
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(
		`SELECT 1 FROM information_schema.constraint_column_usage WHERE constraint_name = $1 AND table_name = $2 AND table_schema = $3`,
		constraintName,
		tableName,
		schemaName,
	).Any()
}

// PredicateIndexExists returns if a index exists on a table in the default schema of the given connection.
func PredicateIndexExists(ctx context.Context, c *db.Connection, tx *sql.Tx, tableName, indexName string) (bool, error) {
	return PredicateIndexExistsInSchema(ctx, c, tx, c.Config.SchemaOrDefault(), tableName, indexName)
}

// PredicateIndexExistsInSchema returns if a index exists on a table in a specific schema on the given connection.
func PredicateIndexExistsInSchema(ctx context.Context, c *db.Connection, tx *sql.Tx, schemaName, tableName, indexName string) (bool, error) {
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(
		`SELECT 1 FROM pg_catalog.pg_indexes where indexname = $1 and tablename = $2 AND schemaname = $3`,
		strings.ToLower(indexName), strings.ToLower(tableName), strings.ToLower(schemaName)).Any()
}

// PredicateRoleExists returns if a role exists or not.
func PredicateRoleExists(ctx context.Context, c *db.Connection, tx *sql.Tx, roleName string) (bool, error) {
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(`SELECT 1 FROM pg_catalog.pg_roles WHERE rolname ilike $1`, roleName).Any()
}

// PredicateSchemaExists returns if a schema exists or not.
func PredicateSchemaExists(ctx context.Context, c *db.Connection, tx *sql.Tx, schemaName string) (bool, error) {
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(
		`SELECT 1 FROM information_schema.schemata WHERE schema_name = $1`,
		schemaName,
	).Any()
}

// PredicateAny returns if a statement has results.
func PredicateAny(ctx context.Context, c *db.Connection, tx *sql.Tx, selectStatement string, params ...interface{}) (bool, error) {
	if !stringutil.HasPrefixCaseless(selectStatement, "select") {
		return false, fmt.Errorf("statement must be a `SELECT`")
	}
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(selectStatement, params...).Any()
}

// PredicateNone returns if a statement doesnt have results.
func PredicateNone(ctx context.Context, c *db.Connection, tx *sql.Tx, selectStatement string, params ...interface{}) (bool, error) {
	if !stringutil.HasPrefixCaseless(selectStatement, "select") {
		return false, fmt.Errorf("statement must be a `SELECT`")
	}
	return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(selectStatement, params...).None()
}

// Not inverts the output of a predicate.
func Not(proceed bool, err error) (bool, error) {
	return !proceed, err
}
