/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stringutil"
)

func TestSerializeState(t *testing.T) {
	assert := assert.New(t)

	state := State{
		RedirectURI: "https://foo.com/bar",
		Token:       stringutil.Random(stringutil.Letters, 32),
		SecureToken: stringutil.Random(stringutil.Letters, 64),
	}

	contents, err := SerializeState(state)
	assert.Nil(err)
	assert.NotEmpty(contents)

	deserialized, err := DeserializeState(contents)
	assert.Nil(err)
	assert.NotNil(deserialized)
	assert.Equal(state.RedirectURI, deserialized.RedirectURI)
	assert.Equal(state.Token, deserialized.Token)
	assert.Equal(state.SecureToken, deserialized.SecureToken)
}
