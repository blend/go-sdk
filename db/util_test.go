package db

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMakeCsvTokens(t *testing.T) {
	a := assert.New(t)

	one := ParamTokensCSV(1)
	two := ParamTokensCSV(2)
	three := ParamTokensCSV(3)

	a.Equal("$1", one)
	a.Equal("$1,$2", two)
	a.Equal("$1,$2,$3", three)
}

func TestReflectSliceType(t *testing.T) {
	assert := assert.New(t)

	objects := []benchObj{
		{}, {}, {},
	}

	ot := ReflectSliceType(objects)
	assert.Equal("benchObj", ot.Name())
}

func TestMakeSliceOfType(t *testing.T) {
	a := assert.New(t)
	tx, txErr := defaultDB().Begin()
	a.Nil(txErr)
	defer func() {
		a.Nil(tx.Rollback())
	}()

	seedErr := seedObjects(10, tx)
	a.Nil(seedErr)

	myType := ReflectType(benchObj{})
	sliceOfT, castOk := makeSliceOfType(myType).(*[]benchObj)
	a.True(castOk)

	allErr := defaultDB().Invoke(OptTx(tx)).All(sliceOfT)
	a.Nil(allErr)
	a.NotEmpty(*sliceOfT)
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

func TestTableName(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("simpletype", TableName(SimpleType{}))
	assert.Equal("simpletype", TableName(&SimpleType{}))
	assert.Equal("not_simple_type_with_name", TableName(SimpleTypeWithName{}))
	assert.Equal("not_simple_type_with_name", TableName(&SimpleTypeWithName{}))

	assert.Equal("embedded_simple_type", TableName(EmbeddedSimpleType{}))
	assert.Equal("embedded_simple_type", TableName(&EmbeddedSimpleType{}))
}
