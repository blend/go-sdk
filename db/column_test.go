package db

import (
	"database/sql"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ref"
)

type jsonFieldValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type setValueTest struct {
	PrimaryKey   string `db:"primary_key,pk"`
	UniqueKey    string `db:"unique_key,uk"`
	InferredName string
	NullFloat64  sql.NullFloat64        `db:"null_float64"`
	NullInt64    sql.NullInt64          `db:"null_int64"`
	NullString   sql.NullString         `db:"null_string"`
	JSON         map[string]interface{} `db:"json_col,json"`
	JSONPtr      *jsonFieldValue        `db:"json_ptr,json"`
	Int64        int64                  `db:"int64"`
	Int64Ptr     *int64                 `db:"int64_ptr"`
	StringPtr    *string                `db:"string_ptr"`
}

func TestSetValue(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)
	col := meta.Lookup()["int64"]
	a.NotNil(col)

	var value int64 = 10
	a.Zero(obj.Int64)
	a.Nil(col.SetValue(&obj, value))
	a.Equal(10, obj.Int64)
}

func TestSetValueConverted(t *testing.T) {
	a := assert.New(t)

	obj := setValueTest{InferredName: "Hello."}
	meta := CachedColumnCollectionFromInstance(obj)
	col := meta.Lookup()["int64"]
	a.NotNil(col)

	value := int(21)
	err := col.SetValue(&obj, value)
	a.Nil(err)
	a.Equal(21, obj.Int64)
}

func TestSetValuePtrAddr(t *testing.T) {
	/*
		Setting a value to an invalid value source shouldn't panic.
	*/
	t.Skip()
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)
	col := meta.Lookup()["string_ptr"]
	a.NotNil(col)

	value := "foobar"
	err := col.SetValue(&obj, value)
	a.NotNil(err)
	a.Nil(obj.StringPtr)
}

func TestSetValueStringPtr(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)
	col := meta.Lookup()["string_ptr"]
	a.NotNil(col)

	value := "foobar"
	err := col.SetValue(&obj, &value)
	a.Nil(err)
	a.NotNil(obj.StringPtr)
	a.Equal("foobar", *obj.StringPtr)
}

func TestSetValueInt64Ptr(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["int64_ptr"]
	a.NotNil(col)
	myValue := int64(21)
	err := col.SetValue(&obj, &myValue)
	a.Nil(err)

	a.NotNil(obj.Int64Ptr)
	a.Equal(21, *obj.Int64Ptr)
}

func TestSetValueJSONNullString(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	err := col.SetValue(&obj,
		sql.NullString{String: `{"foo":"bar"}`, Valid: true},
	)
	a.Nil(err)
	a.Equal("bar", obj.JSON["foo"])
}

func TestSetValueJSONNullStringUnset(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	err := col.SetValue(&obj,
		sql.NullString{},
	)
	a.Nil(err)
	a.Nil(obj.JSON)
}

func TestSetValueJSONNullStringPtr(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	err := col.SetValue(&obj,
		&sql.NullString{String: `{"foo":"bar"}`, Valid: true},
	)
	a.Nil(err)
	a.Equal("bar", obj.JSON["foo"])
}

func TestSetValueJSONNullStringPtrObjectPtr(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_ptr"]
	a.NotNil(col)
	err := col.SetValue(&obj,
		&sql.NullString{String: `{"label":"foo", "value":"bar"}`, Valid: true},
	)
	a.Nil(err)
	a.Equal("foo", obj.JSONPtr.Label)
	a.Equal("bar", obj.JSONPtr.Value)
}

func TestSetValueJSONNullStringPtrUnset(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	err := col.SetValue(&obj,
		&sql.NullString{},
	)
	a.Nil(err)
	a.Nil(obj.JSON)
}

func TestSetValueJSONString(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	err := col.SetValue(&obj,
		`{"foo":"bar"}`,
	)
	a.Nil(err)
	a.Equal("bar", obj.JSON["foo"])
}

func TestSetValueJSONStringPtr(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	value := `{"foo":"bar"}`
	err := col.SetValue(&obj,
		&value,
	)
	a.Nil(err)
	a.Equal("bar", obj.JSON["foo"])
}

func TestSetValueJSONBytes(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	value := []byte(`{"foo":"bar"}`)
	err := col.SetValue(&obj,
		value,
	)
	a.Nil(err)
	a.Equal("bar", obj.JSON["foo"])
}

func TestSetValueJSONBytesPtr(t *testing.T) {
	a := assert.New(t)

	var obj setValueTest
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["json_col"]
	a.NotNil(col)
	value := []byte(`{"foo":"bar"}`)
	err := col.SetValue(&obj,
		&value,
	)
	a.Nil(err)
	a.Equal("bar", obj.JSON["foo"])
}

func TestSetValueResetsNil(t *testing.T) {
	a := assert.New(t)

	obj := setValueTest{Int64Ptr: ref.Int64(1234)}
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["int64_ptr"]
	a.NotNil(col)
	err := col.SetValue(&obj, nil)
	a.Nil(err)
	a.Nil(obj.Int64Ptr)
}

func TestSetValueResetsNilUnset(t *testing.T) {
	a := assert.New(t)

	obj := setValueTest{Int64Ptr: ref.Int64(1234)}
	meta := CachedColumnCollectionFromInstance(obj)

	col := meta.Lookup()["int64_ptr"]
	a.NotNil(col)

	var myValue *int
	err := col.SetValue(&obj, myValue)
	a.Nil(err)
	a.Nil(obj.Int64Ptr)
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
