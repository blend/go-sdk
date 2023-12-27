/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
