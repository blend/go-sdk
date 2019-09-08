package bindata

import "time"

// File is both the file metadata and the contents.
type File struct {
	Name     string
	Contents []byte
	MD5      []byte
	Modtime  time.Time
}
