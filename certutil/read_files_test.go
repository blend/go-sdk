/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestReadFiles(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	files, err := ReadFiles("testdata/client.cert.pem", "testdata/client.key.pem")
	assert.Nil(err)
	assert.Len(files, 2)
}
