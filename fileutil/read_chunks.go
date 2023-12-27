/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package fileutil

import (
	"io"
	"os"

	"github.com/blend/go-sdk/ex"
)

// ReadChunks reads a file in `chunkSize` pieces, dispatched to the handler.
func ReadChunks(filePath string, chunkSize int, handler func([]byte) error) error {
	f, err := os.Open(filePath)
	if err != nil {
		return ex.New(err)
	}
	defer f.Close()

	chunk := make([]byte, chunkSize)
	for {
		readBytes, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		readData := chunk[:readBytes]
		err = handler(readData)
		if err != nil {
			return ex.New(err)
		}
	}
	return nil
}
