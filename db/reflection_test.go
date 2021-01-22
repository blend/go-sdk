package db

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestReflectSliceType(t *testing.T) {
	its := assert.New(t)

	objects := []benchObj{
		{}, {}, {},
	}

	ot := ReflectSliceType(objects)
	its.Equal("benchObj", ot.Name())
}

func TestMakeSliceOfType(t *testing.T) {
	its := assert.New(t)
	tx, txErr := defaultDB().Begin()
	its.Nil(txErr)
	defer func() {
		its.Nil(tx.Rollback())
	}()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	myType := ReflectType(benchObj{})
	sliceOfT, castOk := makeSliceOfType(myType).(*[]benchObj)
	its.True(castOk)

	allErr := defaultDB().Invoke(OptTx(tx)).All(sliceOfT)
	its.Nil(allErr)
	its.NotEmpty(*sliceOfT)
}
