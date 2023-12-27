/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Columns(t *testing.T) {
	its := assert.New(t)

	var emptyColumnCollection ColumnCollection
	firstOrDefaultNil := emptyColumnCollection.FirstOrDefault()
	its.Nil(firstOrDefaultNil)

	obj := myStruct{}
	meta := Columns(obj)

	its.NotNil(meta.Columns())
	its.NotEmpty(meta.Columns())

	its.Equal(9, meta.Len())

	readOnlyColumns := meta.ReadOnly()
	its.Len(readOnlyColumns.Columns(), 1)

	firstOrDefault := meta.FirstOrDefault()
	its.NotNil(firstOrDefault)

	firstCol := meta.FirstOrDefault()
	its.Equal("my_struct", firstCol.TableName)
	its.Equal("PrimaryKeyCol", firstCol.FieldName)
	its.Equal("primary_key_column", firstCol.ColumnName)
	its.True(firstCol.IsPrimaryKey)
	its.True(firstCol.IsAuto)
	its.False(firstCol.IsReadOnly)

	secondCol := meta.Columns()[1]
	its.Equal("auto_column", secondCol.ColumnName)
	its.False(secondCol.IsPrimaryKey)
	its.True(secondCol.IsAuto)
	its.False(secondCol.IsReadOnly)

	thirdCol := meta.Columns()[2]
	its.Equal("InferredName", thirdCol.ColumnName)
	its.False(thirdCol.IsPrimaryKey)
	its.False(thirdCol.IsAuto)
	its.False(thirdCol.IsReadOnly)

	fourthCol := meta.Columns()[3]
	its.Equal("Unique", fourthCol.ColumnName)
	its.False(fourthCol.IsPrimaryKey)
	its.True(fourthCol.IsUniqueKey)
	its.False(fourthCol.IsAuto)
	its.False(fourthCol.IsReadOnly)

	fifthCol := meta.Columns()[4]
	its.Equal("nullable", fifthCol.ColumnName)
	its.False(fifthCol.IsPrimaryKey)
	its.False(fifthCol.IsAuto)
	its.False(fifthCol.IsReadOnly)

	sixthCol := meta.Columns()[5]
	its.Equal("InferredWithFlags", sixthCol.ColumnName)
	its.False(sixthCol.IsPrimaryKey)
	its.False(sixthCol.IsAuto)
	its.True(sixthCol.IsReadOnly)

	uks := meta.UniqueKeys()
	its.Equal(1, uks.Len())
}

func Test_ColumnNameCSV(t *testing.T) {
	its := assert.New(t)

	expected := "primary_key_column,auto_column,InferredName,Unique,nullable,InferredWithFlags,big_int,pointer_col,json_col"
	actual := ColumnNamesCSV(myStruct{})
	its.Equal(expected, actual)
}

func Test_ColumnCollection_Copy(t *testing.T) {
	its := assert.New(t)

	obj := myStruct{}
	meta := Columns(obj)
	newMeta := meta.Copy()
	its.False(meta == newMeta, "These pointers should not be the same.")
	newMeta.columnPrefix = "foo_"
	its.NotEqual(meta.columnPrefix, newMeta.columnPrefix)
}

func Test_ColumnCollection_CopyWithColumnPrefix(t *testing.T) {
	its := assert.New(t)

	obj := myStruct{}
	meta := Columns(obj)
	newMeta := meta.CopyWithColumnPrefix("foo_")
	its.Equal("foo_", newMeta.columnPrefix)
	its.False(meta == newMeta, "These pointers should not be the same.")
	its.NotEqual(meta.columnPrefix, newMeta.columnPrefix)
}

func Test_ColumnCollection_Add(t *testing.T) {
	its := assert.New(t)

	obj := myStruct{}
	meta := Columns(obj)
	newMeta := meta.Copy()
	its.Len(newMeta.columns, 9)
	its.False(newMeta.HasColumn("testo"))
	newMeta.Add(Column{
		FieldName:	"Testo",
		ColumnName:	"testo",
	})
	its.Len(newMeta.columns, 10)
	its.True(newMeta.HasColumn("testo"))
}

func Test_ColumnCollection_Remove(t *testing.T) {
	its := assert.New(t)

	obj := myStruct{}
	meta := Columns(obj)
	newMeta := meta.Copy()

	its.True(newMeta.HasColumn("primary_key_column"))
	newMeta.Remove("primary_key_column")
	its.False(newMeta.HasColumn("primary_key_column"))
}

func Test_ColumnCollection_InsertColumns(t *testing.T) {
	its := assert.New(t)

	obj := myStruct{}
	meta := Columns(obj)
	writeCols := meta.InsertColumns()
	its.NotZero(writeCols.Len())
}

func Test_ColumnCollection_Zero(t *testing.T) {
	its := assert.New(t)
	cols := Columns(columnsZeroTest{})

	emptyCols := cols.Zero(columnsZeroTest{})
	its.Equal(len(cols.columns), len(emptyCols.columns))

	stringValue := "test"
	bytesValue := []byte(stringValue)
	ts := time.Now().UTC()
	setCols := cols.Zero(columnsZeroTest{
		Int:		1,
		Float64:	1,
		TimePtr:	&ts,
		StringPtr:	&stringValue,
		Bytes:		bytesValue,
	})
	its.Equal(2, len(setCols.columns))

	its.False(setCols.HasColumn("Int"))
	its.False(setCols.HasColumn("Float64"))
	its.False(setCols.HasColumn("TimePtr"))
	its.False(setCols.HasColumn("StringPtr"))
	its.False(setCols.HasColumn("Bytes"))

	its.True(setCols.HasColumn("Time"), setCols.String())
	its.True(setCols.HasColumn("String"), setCols.String())

	allCols := cols.Zero(columnsZeroTest{
		Int:		1,
		Float64:	1,
		Time:		ts,
		TimePtr:	&ts,
		String:		stringValue,
		StringPtr:	&stringValue,
		Bytes:		bytesValue,
	})
	its.Empty(allCols.columns)
}

func Test_ColumnCollection_NotZero(t *testing.T) {
	its := assert.New(t)
	cols := Columns(columnsZeroTest{})

	emptyCols := cols.NotZero(columnsZeroTest{})
	its.Empty(emptyCols.columns)

	stringValue := "test"
	bytesValue := []byte(stringValue)
	ts := time.Now().UTC()
	setCols := cols.NotZero(columnsZeroTest{
		Int:		1,
		Float64:	1,
		TimePtr:	&ts,
		StringPtr:	&stringValue,
		Bytes:		bytesValue,
	})
	its.Equal(5, len(setCols.columns))

	its.True(setCols.HasColumn("Int"))
	its.True(setCols.HasColumn("Float64"))
	its.True(setCols.HasColumn("TimePtr"))
	its.True(setCols.HasColumn("StringPtr"))
	its.True(setCols.HasColumn("Bytes"))

	its.False(setCols.HasColumn("Time"), setCols.String())
	its.False(setCols.HasColumn("String"), setCols.String())

	allCols := cols.NotZero(columnsZeroTest{
		Int:		1,
		Float64:	1,
		Time:		ts,
		TimePtr:	&ts,
		String:		stringValue,
		StringPtr:	&stringValue,
		Bytes:		bytesValue,
	})
	its.Len(allCols.columns, 7)
}

func Test_ColumnCollection_ColumnNamesCSVFromAlias(t *testing.T) {
	its := assert.New(t)

	columns := []Column{
		{ColumnName: "foo0"},
		{ColumnName: "foo1"},
		{ColumnName: "foo2"},
	}

	withoutPrefix := &ColumnCollection{
		columns: columns,
	}

	its.Equal("buzz.foo0,buzz.foo1,buzz.foo2", withoutPrefix.ColumnNamesCSVFromAlias("buzz"))

	withPrefix := &ColumnCollection{
		columns:	columns,
		columnPrefix:	"bar_",
	}

	its.Equal("buzz.foo0 as bar_foo0,buzz.foo1 as bar_foo1,buzz.foo2 as bar_foo2", withPrefix.ColumnNamesCSVFromAlias("buzz"))
}

func Test_newColumnCacheKey(t *testing.T) {
	its := assert.New(t)

	its.Equal("db.cacheKeyEmpty", newColumnCacheKey(reflect.TypeOf(cacheKeyEmpty{})))
	its.Equal("db.cacheKeyWithTableName_with_table_name", newColumnCacheKey(reflect.TypeOf(cacheKeyWithTableName{})))
	its.Equal("db.cacheKeyWithColumMetaCacheKeyProvider_with_column_meta_cache_key", newColumnCacheKey(reflect.TypeOf(cacheKeyWithColumMetaCacheKeyProvider{})))
}

//
// helper types
//

type cacheKeyEmpty struct{}

type cacheKeyWithTableName struct{}

func (j cacheKeyWithTableName) TableName() string	{ return "with_table_name" }

type cacheKeyWithColumMetaCacheKeyProvider struct{}

func (j cacheKeyWithColumMetaCacheKeyProvider) ColumnMetaCacheKey() string {
	return "with_column_meta_cache_key"
}

type columnsZeroTest struct {
	Int		int
	Float64		float64
	String		string
	StringPtr	*string
	Time		time.Time
	TimePtr		*time.Time
	Bytes		[]byte
}

type EmbeddedMeta struct {
	PrimaryKeyCol	int	`json:"pk" db:"primary_key_column,pk,serial"`
	AutoCol		string	`json:"auto" db:"auto_column,auto"`
}

type subStruct struct {
	Foo string `json:"foo"`
}

type myStruct struct {
	EmbeddedMeta		`db:",inline"`
	InferredName		string		`json:"normal"`
	Unique			string		`db:",uk"`
	Excluded		string		`json:"-" db:"-"`
	NullableCol		string		`json:"not_nullable" db:"nullable,nullable"`
	InferredWithFlags	string		`db:",readonly"`
	BigIntColumn		int64		`db:"big_int"`
	PointerColumn		*int		`db:"pointer_col"`
	JSONColumn		subStruct	`db:"json_col,json"`
}

func (m myStruct) TableName() string {
	return "my_struct"
}
