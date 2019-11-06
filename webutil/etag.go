package webutil

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
