/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

func Test_Column_SetValue(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)
	col := meta.Lookup()["int64"]
	its.NotNil(col)

	var value int64 = 10
	its.Zero(obj.Int64)
	its.Nil(col.SetValue(&obj, value))
	its.Equal(10, obj.Int64)
}

func Test_Column_SetValueConverted(t *testing.T) {
	its := assert.New(t)

	obj := setValueTest{InferredName: "Hello."}
	meta := Columns(obj)
	col := meta.Lookup()["int64"]
	its.NotNil(col)

	value := int(21)
	err := col.SetValue(&obj, value)
	its.Nil(err)
	its.Equal(21, obj.Int64)
}

func Test_Column_SetValuePtrAddr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)
	col := meta.Lookup()["string_ptr"]
	its.NotNil(col)

	value := "foobar"
	err := col.SetValue(&obj, value)
	its.NotNil(err)
	its.Nil(obj.StringPtr)
}

func Test_Column_SetValueStringPtr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)
	col := meta.Lookup()["string_ptr"]
	its.NotNil(col)

	value := "foobar"
	err := col.SetValue(&obj, &value)
	its.Nil(err)
	its.NotNil(obj.StringPtr)
	its.Equal("foobar", *obj.StringPtr)
}

func Test_Column_SetValueInt64Ptr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["int64_ptr"]
	its.NotNil(col)
	myValue := int64(21)
	err := col.SetValue(&obj, &myValue)
	its.Nil(err)

	its.NotNil(obj.Int64Ptr)
	its.Equal(21, *obj.Int64Ptr)
}

func Test_Column_SetValueJSONNullString(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	err := col.SetValue(&obj,
		sql.NullString{String: `{"foo":"bar"}`, Valid: true},
	)
	its.Nil(err)
	its.Equal("bar", obj.JSON["foo"])
}

func Test_Column_SetValueJSONNullStringUnset(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	err := col.SetValue(&obj,
		sql.NullString{},
	)
	its.Nil(err)
	its.Nil(obj.JSON)
}

func Test_Column_SetValueJSONNullStringPtr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	err := col.SetValue(&obj,
		&sql.NullString{String: `{"foo":"bar"}`, Valid: true},
	)
	its.Nil(err)
	its.Equal("bar", obj.JSON["foo"])
}

func Test_Column_SetValueJSONNullStringPtrObjectPtr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_ptr"]
	its.NotNil(col)
	err := col.SetValue(&obj,
		&sql.NullString{String: `{"label":"foo", "value":"bar"}`, Valid: true},
	)
	its.Nil(err)
	its.Equal("foo", obj.JSONPtr.Label)
	its.Equal("bar", obj.JSONPtr.Value)
}

func Test_Column_SetValueJSONNullStringPtrUnset(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	err := col.SetValue(&obj,
		&sql.NullString{},
	)
	its.Nil(err)
	its.Nil(obj.JSON)
}

func Test_Column_SetValueJSONString(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	err := col.SetValue(&obj,
		`{"foo":"bar"}`,
	)
	its.Nil(err)
	its.Equal("bar", obj.JSON["foo"])
}

func Test_Column_SetValueJSONStringPtr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	value := `{"foo":"bar"}`
	err := col.SetValue(&obj,
		&value,
	)
	its.Nil(err)
	its.Equal("bar", obj.JSON["foo"])
}

func Test_Column_SetValueJSONBytes(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	value := []byte(`{"foo":"bar"}`)
	err := col.SetValue(&obj,
		value,
	)
	its.Nil(err)
	its.Equal("bar", obj.JSON["foo"])
}

func Test_Column_SetValueJSONBytesPtr(t *testing.T) {
	its := assert.New(t)

	var obj setValueTest
	meta := Columns(obj)

	col := meta.Lookup()["json_col"]
	its.NotNil(col)
	value := []byte(`{"foo":"bar"}`)
	err := col.SetValue(&obj,
		&value,
	)
	its.Nil(err)
	its.Equal("bar", obj.JSON["foo"])
}

func Test_Column_SetValueResetsNil(t *testing.T) {
	its := assert.New(t)

	obj := setValueTest{Int64Ptr: ref.Int64(1234)}
	meta := Columns(obj)

	col := meta.Lookup()["int64_ptr"]
	its.NotNil(col)
	err := col.SetValue(&obj, nil)
	its.Nil(err)
	its.Nil(obj.Int64Ptr)
}

func Test_Column_SetValueResetsNilUnset(t *testing.T) {
	its := assert.New(t)

	obj := setValueTest{Int64Ptr: ref.Int64(1234)}
	meta := Columns(obj)

	col := meta.Lookup()["int64_ptr"]
	its.NotNil(col)

	var myValue *int
	err := col.SetValue(&obj, myValue)
	its.Nil(err)
	its.Nil(obj.Int64Ptr)
}

func Test_Column_GetValue(t *testing.T) {
	its := assert.New(t)
	obj := myStruct{EmbeddedMeta: EmbeddedMeta{PrimaryKeyCol: 5}, InferredName: "Hello."}
	meta := Columns(obj)
	pk := meta.PrimaryKeys().FirstOrDefault()
	its.NotNil(pk)
	value := pk.GetValue(&obj)
	its.NotNil(value)
	its.Equal(5, value)
}
