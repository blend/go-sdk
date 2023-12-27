/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/blend/go-sdk/ex"
)

// ETag creates an etag for a given blob.
func ETag(contents []byte) (string, error) {
	hash := md5.New()
	_, err := hash.Write(contents)
	if err != nil {
		return "", ex.New(err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
