package webutil

import "net/http"

// GetContentType gets the content type out of a header collection.
func GetContentType(header http.Header) string {
	if header != nil {
		header.Get(HeaderContentType)
	}
	return ""
}
