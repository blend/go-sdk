/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package crypto

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()

	its := assert.New(t)
	key, err := CreateKey(32)
	its.Nil(err)
	plaintext := "Mary Jane Hawkins"

	ciphertext, err := Encrypt(key, []byte(plaintext))
	its.Nil(err)

	decryptedPlaintext, err := Decrypt(key, ciphertext)
	its.Nil(err)
	its.Equal(plaintext, string(decryptedPlaintext))
}
