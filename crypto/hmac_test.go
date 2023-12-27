/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package crypto

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_HMAC(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	key, err := CreateKey(128)
	its.Nil(err)
	plaintext := "123-12-1234"
	its.Equal(
		HMAC512(key, []byte(plaintext)),
		HMAC512(key, []byte(plaintext)),
	)
}
