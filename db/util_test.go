/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"reflect"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_ParamTokensCSV(t *testing.T) {
	its := assert.New(t)

	one := ParamTokensCSV(1)
	two := ParamTokensCSV(2)
	three := ParamTokensCSV(3)

	its.Equal("$1", one)
	its.Equal("$1,$2", two)
	its.Equal("$1,$2,$3", three)
}

type SimpleType struct {
	ID   int
	Name string
}

type SimpleTypeWithName struct {
	ID   int
	Name string
}

func (st SimpleTypeWithName) TableName() string {
	return "not_simple_type_with_name"
}

type EmbeddedSimpleTypeMeta struct {
	ID   int
	Name string
}

type EmbeddedSimpleType struct {
	EmbeddedSimpleTypeMeta
}

func (est EmbeddedSimpleType) TableName() string {
	return "embedded_simple_type"
}

func Test_TableName(t *testing.T) {
	its := assert.New(t)

	its.Equal("SimpleType", TableName(SimpleType{}))
	its.Equal("SimpleType", TableName(new(SimpleType)))

	its.Equal("not_simple_type_with_name", TableName(SimpleTypeWithName{}))
	its.Equal("not_simple_type_with_name", TableName(new(SimpleTypeWithName)))

	its.Equal("embedded_simple_type", TableName(EmbeddedSimpleType{}))
	its.Equal("embedded_simple_type", TableName(new(EmbeddedSimpleType)))
}

func Test_TableNameByType(t *testing.T) {
	its := assert.New(t)

	tt := reflect.TypeOf(SimpleType{})
	its.Equal("SimpleType", TableNameByType(tt))

	tt = reflect.TypeOf(new(SimpleType))
	its.Equal("SimpleType", TableNameByType(tt))

	tt = reflect.TypeOf(SimpleTypeWithName{})
	its.Equal("not_simple_type_with_name", TableNameByType(tt))

	tt = reflect.TypeOf(new(SimpleTypeWithName))
	its.Equal("not_simple_type_with_name", TableNameByType(tt))

	tt = reflect.TypeOf(EmbeddedSimpleType{})
	its.Equal("embedded_simple_type", TableNameByType(tt))

	tt = reflect.TypeOf(new(EmbeddedSimpleType))
	its.Equal("embedded_simple_type", TableNameByType(tt))
}
