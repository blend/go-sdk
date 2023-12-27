/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

// DatabaseMapped is the interface that any objects passed into database mapped methods like Create, Update, Delete, Get, GetAll etc.
type DatabaseMapped interface{}

// TableNameProvider is a type that implements the TableName() function.
// The only required method is TableName() string that returns the name of the table in the database this type is mapped to.
//
//	type MyDatabaseMappedObject {
//		Mycolumn `db:"my_column"`
//	}
//	func (_ MyDatabaseMappedObject) TableName() string {
//		return "my_database_mapped_object"
//	}
//
// If you require different table names based on alias, create another type.
type TableNameProvider interface {
	TableName() string
}

// ColumnMetaCacheKeyProvider is a provider for a column meta key.
type ColumnMetaCacheKeyProvider interface {
	ColumnMetaCacheKey() string
}

// Populatable is an interface that you can implement if your object is read often and is performance critical.
type Populatable interface {
	Populate(rows Scanner) error
}

// RowsConsumer is the function signature that is called from within Each().
type RowsConsumer func(r Rows) error

// Scanner is a type that can scan into variadic values.
type Scanner interface {
	Scan(...interface{}) error
}

// ColumnsProvider is a type that can return columns.
type ColumnsProvider interface {
	Columns() ([]string, error)
}

// Rows provides the relevant fields to populate by name.
type Rows interface {
	Scanner
	ColumnsProvider
}
