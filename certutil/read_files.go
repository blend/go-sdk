package certutil

import (
	"io/ioutil"
	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrInvalidCertPEM ex.Class = "failed to add cert to pool as pem"
)

// MustBytes panics on an error or returns the contents.
func MustBytes(contents []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return contents
}

// BytesWithError returns a bytes error response with the error
// as an ex.
func BytesWithError(bytes []byte, err error) ([]byte, error) {
	return bytes, ex.New(err)
}

// ReadFiles reads a list of files as bytes.
func ReadFiles(files ...string) (data [][]byte, err error) {
	var contents []byte
	for _, path := range files {
		contents, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, ex.New(err)
		}
		data = append(data, contents)
	}
	return
}
