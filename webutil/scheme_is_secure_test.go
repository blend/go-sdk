/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_SchemeIsSecure(t *testing.T) {
	its := assert.New(t)

	its.True(SchemeIsSecure(SchemeHTTPS))
	its.True(SchemeIsSecure(SchemeSPDY))

	its.False(SchemeIsSecure(SchemeHTTP))
	its.False(SchemeIsSecure("garbage"))
	its.False(SchemeIsSecure(""))
}
