/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
