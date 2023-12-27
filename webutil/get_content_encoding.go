/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import "net/http"

// GetContentEncoding gets the content type out of a header collection.
func GetContentEncoding(header http.Header) string {
	if header != nil {
		return header.Get(HeaderContentEncoding)
	}
	return ""
}
