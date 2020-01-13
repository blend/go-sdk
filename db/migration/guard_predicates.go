package migration

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/db"
)

// TableExists returns a guard that ensures a table exists
func TableExists(tableName string) GuardFunc {
	return guardPredicate(fmt.Sprintf("Check table exists: %s", tableName), PredicateTableExists, tableName)
}

// TableNotExists returns a guard that ensures a table does not exist
func TableNotExists(tableName string) GuardFunc {
	return guardNotPredicate(fmt.Sprintf("Check table does not exist: %s", tableName), PredicateTableExists, tableName)
}

// TableExistsInSchema returns a guard that ensures a table exists
func TableExistsInSchema(schemaName, tableName string) GuardFunc {
	return guardPredicate2(fmt.Sprintf("Check table exists: %s.%s", schemaName, tableName),
		PredicateTableExistsInSchema, schemaName, tableName)
}

// TableNotExistsInSchema returns a guard that ensures a table exists
func TableNotExistsInSchema(schemaName, tableName string) GuardFunc {
	return guardNotPredicate2(fmt.Sprintf("Check table does not exist: %s.%s", schemaName, tableName),
		PredicateTableExistsInSchema, schemaName, tableName)
}

// ColumnExists returns a guard that ensures a column exists
func ColumnExists(tableName, columnName string) GuardFunc {
	return guardPredicate2(fmt.Sprintf("Check column exists: %s.%s", tableName, columnName),
		PredicateColumnExists, tableName, columnName)
}

// ColumnNotExists returns a guard that ensures a column does not exist
func ColumnNotExists(tableName, columnName string) GuardFunc {
	return guardNotPredicate2(fmt.Sprintf("Check column does not exist: %s.%s", tableName, columnName),
		PredicateColumnExists, tableName, columnName)
}

// ColumnExistsInSchema returns a guard that ensures a column exists
func ColumnExistsInSchema(schemaName, tableName, columnName string) GuardFunc {
	return guardPredicate3(fmt.Sprintf("Check column exists: %s.%s.%s", schemaName, tableName, columnName),
		PredicateColumnExistsInSchema, schemaName, tableName, columnName)
}

// ColumnNotExistsInSchema returns a guard that ensures a column does not exist
func ColumnNotExistsInSchema(schemaName, tableName, columnName string) GuardFunc {
	return guardNotPredicate3(fmt.Sprintf("Check column does not exist: %s.%s.%s", schemaName, tableName, columnName),
		PredicateColumnExistsInSchema, schemaName, tableName, columnName)
}

// ConstraintExists returns a guard that ensures a constraint exists
func ConstraintExists(tableName, constraintName string) GuardFunc {
	return guardPredicate2(fmt.Sprintf("Check constraint %s exists on table %s", constraintName, tableName),
		PredicateConstraintExists, tableName, constraintName)
}

// ConstraintNotExists returns a guard that ensures a constraint does not exist
func ConstraintNotExists(tableName, constraintName string) GuardFunc {
	return guardNotPredicate2(fmt.Sprintf("Check constraint %s does not exist on table %s", constraintName, tableName),
		PredicateConstraintExists, tableName, constraintName)
}

// ConstraintExistsInSchema returns a guard that ensures a constraint exists
func ConstraintExistsInSchema(schemaName, tableName, constraintName string) GuardFunc {
	return guardPredicate3(fmt.Sprintf("Check constraint %s exists on table %s.%s", constraintName, schemaName, tableName),
		PredicateConstraintExistsInSchema, schemaName, tableName, constraintName)
}

// ConstraintNotExistsInSchema returns a guard that ensures a constraint does not exist
func ConstraintNotExistsInSchema(schemaName, tableName, constraintName string) GuardFunc {
	return guardNotPredicate3(fmt.Sprintf("Check constraint %s does not exist on table %s.%s", constraintName, schemaName, tableName),
		PredicateConstraintExistsInSchema, schemaName, tableName, constraintName)
}

// IndexExists returns a guard that ensures an index exists
func IndexExists(tableName, indexName string) GuardFunc {
	return guardPredicate2(fmt.Sprintf("Check index %s exists on table %s", indexName, tableName),
		PredicateIndexExists, tableName, indexName)
}

// IndexNotExists returns a guard that ensures an index does not exist
func IndexNotExists(tableName, indexName string) GuardFunc {
	return guardNotPredicate2(fmt.Sprintf("Check index %s does not exist on table %s", indexName, tableName),
		PredicateIndexExists, tableName, indexName)
}

// IndexExistsInSchema returns a guard that ensures an index exists
func IndexExistsInSchema(schemaName, tableName, indexName string) GuardFunc {
	return guardPredicate3(fmt.Sprintf("Check index %s exists on table %s.%s", indexName, schemaName, tableName),
		PredicateIndexExistsInSchema, schemaName, tableName, indexName)
}

// IndexNotExistsInSchema returns a guard that ensures an index does not exist
func IndexNotExistsInSchema(schemaName, tableName, indexName string) GuardFunc {
	return guardNotPredicate3(fmt.Sprintf("Check index %s does not exist on table %s.%s", indexName, schemaName, tableName),
		PredicateIndexExistsInSchema, schemaName, tableName, indexName)
}

// RoleExists returns a guard that ensures a role (user) exists
func RoleExists(roleName string) GuardFunc {
	return guardPredicate(fmt.Sprintf("Check Role Exists: %s", roleName), PredicateRoleExists, roleName)
}

// RoleNotExists returns a guard that ensures a role (user) does not exist
func RoleNotExists(roleName string) GuardFunc {
	return guardNotPredicate(fmt.Sprintf("Check Role Not Exists: %s", roleName), PredicateRoleExists, roleName)
}

// SchemaExists is a guard function for asserting that a schema exists
func SchemaExists(schemaName string) GuardFunc {
	return Guard(fmt.Sprintf("drop schema `%s`", schemaName),
		func(c *db.Connection, tx *sql.Tx) (bool, error) {
			return predicateSchemaExists(c, tx, schemaName)
		})
}

// SchemaNotExists is a guard function for asserting that a schema does not exist
func SchemaNotExists(schemaName string) GuardFunc {
	return Guard(fmt.Sprintf("create schema `%s`", schemaName),
		func(c *db.Connection, tx *sql.Tx) (bool, error) {
			return Not(predicateSchemaExists(c, tx, schemaName))
		})
}

func predicateSchemaExists(c *db.Connection, tx *sql.Tx, schemaName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM information_schema.schemata WHERE schema_name = $1`,
		strings.ToLower(schemaName)).Any()
}
