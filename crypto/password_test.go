/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package crypto

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_PasswordHashAndMatch(t *testing.T) {
	t.Parallel()
	its := assert.New(t)
	password := "some-test-password-12345"
	hashedPassword, err := HashPassword(password)
	its.Nil(err)
	its.NotEqual("", hashedPassword)
	its.True(PasswordMatchesHash(password, hashedPassword))
	its.False(PasswordMatchesHash("something-else", hashedPassword))
}
