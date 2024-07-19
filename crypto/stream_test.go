/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"bytes"
	"io"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Stream_EncrypterDecrypter(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	encKey, err := CreateKey(32)
	its.Nil(err)
	macKey, err := CreateKey(32)
	its.Nil(err)
	plaintext := "Eleven is the best person in all of Hawkins Indiana. Some more text"
	pt := []byte(plaintext)

	src := bytes.NewReader(pt)

	se, err := NewStreamEncrypter(encKey, macKey, src)
	its.Nil(err)
	its.NotNil(se)

	encrypted, err := io.ReadAll(se)
	its.Nil(err)
	its.NotNil(encrypted)

	sd, err := NewStreamDecrypter(encKey, macKey, se.Meta(), bytes.NewReader(encrypted))
	its.Nil(err)
	its.NotNil(sd)

	decrypted, err := io.ReadAll(sd)
	its.Nil(err)
	its.Equal(plaintext, string(decrypted))

	its.Nil(sd.Authenticate())
}

func Test_checkedWrite(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	writer := bytes.NewBuffer(nil)
	data := []byte{1, 1, 1}
	v, err := checkedWrite(writer, data)
	its.Nil(err)
	its.Equal(len(data), v)
}
