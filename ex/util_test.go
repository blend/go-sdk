/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

type classProvider struct {
	error
	ErrClass	error
}

func (cp classProvider) Class() error {
	return cp.ErrClass
}

func TestErrClass(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(ErrClass(nil))
	var unsetErr error
	assert.Nil(ErrClass(unsetErr))

	assert.Nil(ErrClass("foo"))

	err := New("this is a test")
	assert.Equal("this is a test", ErrClass(err).Error())

	cp := classProvider{
		error:		fmt.Errorf("this is a provider test"),
		ErrClass:	fmt.Errorf("the error class"),
	}
	assert.Equal("the error class", ErrClass(cp).Error())
	assert.Equal("this is a test", ErrClass(fmt.Errorf("this is a test")).Error())
}

func TestErrMessage(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(ErrMessage(nil))
	assert.Empty(ErrMessage(fmt.Errorf("foo bar baz")))
	assert.Equal("this is a message", ErrMessage(New("error class", OptMessage("this is a message"))))
}

type stackProvider struct {
	error
	Stack	StackTrace
}

func (sp stackProvider) StackTrace() StackTrace {
	return sp.Stack
}

func TestErrStackTrace(t *testing.T) {
	assert := assert.New(t)

	err := New("this is a test")
	assert.NotNil(ErrStackTrace(err))

	sp := stackProvider{
		error:	fmt.Errorf("this is a provider test"),
		Stack:	StackStrings([]string{"first", "second"}),
	}
	assert.Equal([]string{"first", "second"}, ErrStackTrace(sp).Strings())

	assert.Nil(ErrStackTrace(fmt.Errorf("this is also a test")))
}
