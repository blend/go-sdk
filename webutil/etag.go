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
