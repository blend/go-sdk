package db

import (
	"reflect"
	"testing"
	"github.com/blend/go-sdk/assert"
)

type EmbeddedMeta struct {
	PrimaryKeyCol int    `json:"pk" db:"primary_key_column,pk,serial"`
	AutoCol       string `json:"auto" db:"auto_column,auto"`
}

type subStruct struct {
	Foo string `json:"foo"`
}

type myStruct struct {
	EmbeddedMeta      `db:",inline"`
	InferredName      string    `json:"normal"`
	Unique            string    `db:",uk"`
	Excluded          string    `json:"-" db:"-"`
	NullableCol       string    `json:"not_nullable" db:"nullable,nullable"`
	InferredWithFlags string    `db:",readonly"`
	BigIntColumn      int64     `db:"big_int"`
	PointerColumn     *int      `db:"pointer_col"`
	JSONColumn        subStruct `db:"json_col,json"`
}

func (m myStruct) TableName() string {
	return "my_struct"
}

func TestGetCachedColumnCollectionFromInstance(t *testing.T) {
	a := assert.New(t)

	emptyColumnCollection := ColumnCollection{}
	firstOrDefaultNil := emptyColumnCollection.FirstOrDefault()
	a.Nil(firstOrDefaultNil)

	obj := myStruct{}
	meta := CachedColumnCollectionFromInstance(obj)

	a.NotNil(meta.Columns())
	a.NotEmpty(meta.Columns())

	a.Equal(9, meta.Len())

	readOnlyColumns := meta.ReadOnly()
	a.Len(readOnlyColumns.Columns(), 1)

	firstOrDefault := meta.FirstOrDefault()
	a.NotNil(firstOrDefault)

	firstCol := meta.FirstOrDefault()
	a.Equal("my_struct", firstCol.TableName)
	a.Equal("PrimaryKeyCol", firstCol.FieldName)
	a.Equal("primary_key_column", firstCol.ColumnName)
	a.True(firstCol.IsPrimaryKey)
	a.True(firstCol.IsAuto)
	a.False(firstCol.IsReadOnly)

	secondCol := meta.Columns()[1]
	a.Equal("auto_column", secondCol.ColumnName)
	a.False(secondCol.IsPrimaryKey)
	a.True(secondCol.IsAuto)
	a.False(secondCol.IsReadOnly)

	thirdCol := meta.Columns()[2]
	a.Equal("inferredname", thirdCol.ColumnName)
	a.False(thirdCol.IsPrimaryKey)
	a.False(thirdCol.IsAuto)
	a.False(thirdCol.IsReadOnly)

	fourthCol := meta.Columns()[3]
	a.Equal("unique", fourthCol.ColumnName)
	a.False(fourthCol.IsPrimaryKey)
	a.True(fourthCol.IsUniqueKey)
	a.False(fourthCol.IsAuto)
	a.False(fourthCol.IsReadOnly)

	fifthCol := meta.Columns()[4]
	a.Equal("nullable", fifthCol.ColumnName)
	a.False(fifthCol.IsPrimaryKey)
	a.False(fifthCol.IsAuto)
	a.False(fifthCol.IsReadOnly)

	sixthCol := meta.Columns()[5]
	a.Equal("inferredwithflags", sixthCol.ColumnName)
	a.False(sixthCol.IsPrimaryKey)
	a.False(sixthCol.IsAuto)
	a.True(sixthCol.IsReadOnly)

	uks := meta.UniqueKeys()
	a.Equal(1, uks.Len())
}

func TestColumnCollectionCopy(t *testing.T) {
	assert := assert.New(t)

	obj := myStruct{}
	meta := CachedColumnCollectionFromInstance(obj)
	newMeta := meta.Copy()
	assert.False(meta == newMeta, "These pointers should not be the same.")
	newMeta.columnPrefix = "foo_"
	assert.NotEqual(meta.columnPrefix, newMeta.columnPrefix)
}

func TestColumnCollectionRemove(t *testing.T) {
	assert := assert.New(t)

	obj := myStruct{}
	meta := CachedColumnCollectionFromInstance(obj)
	newMeta := meta.Copy()

	assert.True(newMeta.HasColumn("primary_key_column"))
	newMeta.Remove("primary_key_column")
	assert.False(newMeta.HasColumn("primary_key_column"))
}

func TestColumnCollectionWithColumnPrefix(t *testing.T) {
	assert := assert.New(t)

	obj := myStruct{}
	meta := CachedColumnCollectionFromInstance(obj)
	newMeta := meta.CopyWithColumnPrefix("foo_")
	assert.Equal("foo_", newMeta.columnPrefix)
	assert.False(meta == newMeta, "These pointers should not be the same.")
	assert.NotEqual(meta.columnPrefix, newMeta.columnPrefix)
}

func TestColumnCollectionWriteColumns(t *testing.T) {
	assert := assert.New(t)

	obj := myStruct{}
	meta := CachedColumnCollectionFromInstance(obj)
	writeCols := meta.WriteColumns()
	assert.NotZero(writeCols.Len())
}

type cacheKeyEmpty struct{}

type cacheKeyWithTableName struct{}

func (j cacheKeyWithTableName) TableName() string { return "with_table_name" }

type cacheKeyWithColumMetaCacheKeyProvider struct{}

func (j cacheKeyWithColumMetaCacheKeyProvider) ColumnMetaCacheKey() string {
	return "with_column_meta_cache_key"
}

func TestNewColumnCacheKey(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("db.cacheKeyEmpty", newColumnCacheKey(reflect.TypeOf(cacheKeyEmpty{})))
	assert.Equal("db.cacheKeyWithTableName_with_table_name", newColumnCacheKey(reflect.TypeOf(cacheKeyWithTableName{})))
	assert.Equal("db.cacheKeyWithColumMetaCacheKeyProvider_with_column_meta_cache_key", newColumnCacheKey(reflect.TypeOf(cacheKeyWithColumMetaCacheKeyProvider{})))

}
