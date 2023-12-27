/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_Errors_All(t *testing.T) {
	its := assert.New(t)

	errors := Errors(make(chan error, 5))
	its.Nil(errors.All())
	errors <- nil
	errors <- fmt.Errorf("this is just a test 0")
	errors <- nil
	errors <- fmt.Errorf("this is just a test 1")
	errors <- nil

	allErr := errors.All()
	its.NotNil(allErr)

	typed, ok := allErr.(ex.Multi)
	its.True(ok)
	its.Len(typed, 2)
	its.Equal("this is just a test 0", ex.ErrClass(typed[0]).Error())
	its.Equal("this is just a test 1", ex.ErrClass(typed[1]).Error())
}

func Test_Errors_First(t *testing.T) {
	its := assert.New(t)

	errors := Errors(make(chan error, 5))
	its.Nil(errors.All())
	errors <- nil
	errors <- fmt.Errorf("this is just a test 0")
	errors <- nil
	errors <- fmt.Errorf("this is just a test 1")
	errors <- nil

	firstErr := errors.First()
	its.NotNil(firstErr)

	typed, ok := firstErr.(*ex.Ex)
	its.True(ok)
	its.Equal("this is just a test 0", ex.ErrClass(typed).Error())
}
