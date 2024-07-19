/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestListener(t *testing.T) {
	assert := assert.New(t)

	tcpln, err := CreateListener("127.0.0.1:")
	assert.Nil(err)
	defer func() { _ = tcpln.Close() }()
	assert.Equal("tcp", tcpln.Addr().Network())
	assert.Contains(tcpln.Addr().String(), "127.0.0.1:")

	socketDir := os.TempDir()
	socketPath := filepath.Join(socketDir, uuid.V4().String())
	socketAddress := fmt.Sprintf("unix://" + socketPath)
	unixln, err := CreateListener(socketAddress)
	assert.Nil(err)
	defer func() { _ = unixln.Close() }()
	assert.Equal("unix", unixln.Addr().Network())
	assert.Equal(socketPath, unixln.Addr().String())
}
