/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"crypto/md5"
	"encoding/hex"
)

// ETag creates an etag for a given blob.
func ETag(contents []byte) string {
	hash := md5.New()
	_, _ = hash.Write(contents)
	return hex.EncodeToString(hash.Sum(nil))
}
