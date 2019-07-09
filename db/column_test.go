package db

import (
	"database/sql"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"
)

func TestSetValue(t *testing.T) {
	a := assert.New(t)

	obj := myStruct{InferredName: "Hello."}

	var value interface{}
	value = 10
	meta := CachedColumnCollectionFromInstance(obj)
	pk := meta.Columns()[0]
	a.Nil(pk.SetValue(&obj, value))
	a.Equal(10, obj.PrimaryKeyCol)
}

func TestSetValueConverted(t *testing.T) {
	a := assert.New(t)

	obj := myStruct{InferredName: "Hello."}

	meta := CachedColumnCollectionFromInstance(obj)
	col := meta.Lookup()["big_int"]
	a.NotNil(col)
	err := col.SetValue(&obj, int(21))
	a.Nil(err)
	a.Equal(21, obj.BigIntColumn)
}

func TestSetValueJSON(t *testing.T) {
	a := assert.New(t)

	obj := myStruct{InferredName: "Hello."}
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	err := col.SetValue(&obj, sql.NullString{String: `{"foo":"bar"}`, Valid: true})
	a.Nil(err)
	a.Equal("bar", obj.JSONColumn.Foo)
}

func TestSetValuePtr(t *testing.T) {
	a := assert.New(t)

	obj := myStruct{InferredName: "Hello."}
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["pointer_col"]
	a.NotNil(col)
	myValue := 21
	err := col.SetValue(&obj, &myValue)
	a.Nil(err)
	a.NotNil(obj.PointerColumn)
	a.Equal(21, *obj.PointerColumn)
}

func TestSetValueNil(t *testing.T) {
	a := assert.New(t)

	obj := myStruct{PointerColumn: ref.Int(1234)}
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["pointer_col"]
	a.NotNil(col)

	var myValue *int
	err := col.SetValue(&obj, myValue)
	a.Nil(err)
	a.Nil(obj.PointerColumn)
}

func TestGetValue(t *testing.T) {
	a := assert.New(t)
	obj := myStruct{EmbeddedMeta: EmbeddedMeta{PrimaryKeyCol: 5}, InferredName: "Hello."}
	meta := CachedColumnCollectionFromInstance(obj)
	pk := meta.PrimaryKeys().FirstOrDefault()
	a.NotNil(pk)
	value := pk.GetValue(&obj)
	a.NotNil(value)
	a.Equal(5, value)
}
