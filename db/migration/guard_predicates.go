package migration

// ColumnNotExists creates a table on the given connection if it does not exist.
func ColumnNotExists(tableName, columnName string) GuardFunc {
	return ColumnNotExistsWithPredicate(PredicateColumnExists, tableName, columnName)
}

// ConstraintNotExists creates a table on the given connection if it does not exist.
func ConstraintNotExists(constraintName string) GuardFunc {
	return ConstraintNotExistsWithPredicate(PredicateConstraintExists, constraintName)
}

// TableNotExists creates a table on the given connection if it does not exist.
func TableNotExists(tableName string) GuardFunc {
	return TableNotExistsWithPredicate(PredicateTableExists, tableName)
}

// IndexNotExists creates a index on the given connection if it does not exist.
func IndexNotExists(tableName, indexName string) GuardFunc {
	return IndexNotExistsWithPredicate(PredicateIndexExists, tableName, indexName)
}

// RoleNotExists creates a new role if it doesn't exist.
func RoleNotExists(roleName string) GuardFunc {
	return RoleNotExistsWithPredicate(PredicateRoleExists, roleName)
}

// ColumnExists alters an existing column, erroring if it doesn't exist
func ColumnExists(tableName, columnName string) GuardFunc {
	return ColumnExistsWithPredicate(PredicateColumnExists, tableName, columnName)
}

// ConstraintExists alters an existing constraint, erroring if it doesn't exist
func ConstraintExists(constraintName string) GuardFunc {
	return ConstraintExistsWithPredicate(PredicateConstraintExists, constraintName)
}

// TableExists alters an existing table, erroring if it doesn't exist
func TableExists(tableName string) GuardFunc {
	return TableExistsWithPredicate(PredicateTableExists, tableName)
}

// IndexExists alters an existing index, erroring if it doesn't exist
func IndexExists(tableName, indexName string) GuardFunc {
	return IndexExistsWithPredicate(PredicateIndexExists, tableName, indexName)
}

// RoleExists alters an existing role in the db
func RoleExists(roleName string) GuardFunc {
	return RoleExistsWithPredicate(PredicateRoleExists, roleName)
}
