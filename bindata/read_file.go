package bindata

import (
	"crypto/md5"
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

	md5Hash := md5.New().Sum(contents)

	return &File{
		Name:     path,
		Modtime:  stat.ModTime(),
		MD5:      md5Hash,
		Contents: contents,
	}, nil
}
