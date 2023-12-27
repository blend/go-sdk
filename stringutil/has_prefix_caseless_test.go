/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHasPrefixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.True(HasPrefixCaseless("hello world!", "hello"))
	assert.True(HasPrefixCaseless("hello world", "hello world"))
	assert.True(HasPrefixCaseless("HELLO world", "hello"))
	assert.True(HasPrefixCaseless("hello world", "HELLO"))
	assert.True(HasPrefixCaseless("hello world", "h"))

	assert.False(HasPrefixCaseless("hello world", "butters"))
	assert.False(HasPrefixCaseless("hello world", "hello world boy is this long"))
	assert.False(HasPrefixCaseless("hello world", "world"))	//this would pass suffix
}
