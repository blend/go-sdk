/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestToFileCreate(t *testing.T) {
	assert := assert.New(t)

	// create a new file
	filename := uuid.V4().String() + ".temp"
	defer func() {
		os.Remove(filename)
	}()
	file, err := ToFile(filename)
	assert.Nil(err)
	_, err = file.Stat()
	assert.Nil(err)
	assert.Nil(file.Close())
}

func TestToFileOpen(t *testing.T) {
	assert := assert.New(t)

	// create a new file
	filename := uuid.V4().String() + ".temp"
	assert.Nil(Touch(filename))
	defer func() {
		os.Remove(filename)
	}()

	file, err := ToFile(filename)
	assert.Nil(err)
	_, err = file.Stat()
	assert.Nil(err)
}
