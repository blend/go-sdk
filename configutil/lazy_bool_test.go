/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLazybool(t *testing.T) {
	its := assert.New(t)

	isNil := LazyBool(nil)
	var value bool = false
	hasValue := LazyBool(&value)
	var value2 bool = true
	hasValue2 := LazyBool(&value2)

	var setValue bool
	its.Nil(SetBool(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	its.Equal(true, setValue)
}
