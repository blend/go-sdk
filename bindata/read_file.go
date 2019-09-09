package bindata

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"

	"github.com/blend/go-sdk/ex"
)

// ReadFile reads a file at a given path.
func ReadFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ex.New(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, ex.New(err)
	}

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, ex.New(err)
	}

	hasher := md5.New()
	if _, err := io.Copy(hasher, bytes.NewReader(contents)); err != nil {
		return nil, ex.New(err)
	}

	return &File{
		Name:     path,
		Modtime:  stat.ModTime(),
		MD5:      hasher.Sum(nil),
		Contents: contents,
	}, nil
}
