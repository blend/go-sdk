/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

type testConfig struct {
	Omitted      string          `json:"omitted"`
	ExplicitOmit string          `json:"explicitOmit" secret:"-"`
	Included     string          `json:"included" secret:"included"`
	Number       float64         `json:"number" secret:"number"`
	Binary       []byte          `json:"binary" secret:"binary"`
	Nested       testConfigInner `json:"nested" secret:"nested"`
}

type testConfigInner struct {
	Foo string
	Bar string
}

func TestDecomposeRestore(t *testing.T) {
	assert := assert.New(t)

	config := testConfig{
		Omitted:      "a",
		ExplicitOmit: "b",
		Included:     "c",
		Number:       3.14,
		Binary:       []byte("just a test"),
		Nested: testConfigInner{
			Foo: "is foo",
			Bar: "is bar",
		},
	}

	data, err := DecomposeJSON(config)
	assert.Nil(err)
	assert.Len(data, 4)

	var verify testConfig
	assert.Nil(RestoreJSON(data, &verify))
	assert.Empty(verify.Omitted)
	assert.Empty(verify.ExplicitOmit)
	assert.Equal("c", verify.Included)
	assert.Equal(3.14, verify.Number)
	assert.Equal([]byte("just a test"), verify.Binary)
	assert.Equal("is foo", verify.Nested.Foo)
	assert.Equal("is bar", verify.Nested.Bar)
}
