/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_JSON(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var foo *int

	marshaled, err := json.Marshal(foo)
	its.Nil(err)
	its.Equal("null", string(marshaled))
	its.Nil(JSON(foo))

	valid := struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}{
		Foo: "not-foo",
		Bar: "not-bar",
	}

	output := JSON(valid)
	its.NotNil(output)
	its.Equal(`{"foo":"not-foo","bar":"not-bar"}`, output)
}
