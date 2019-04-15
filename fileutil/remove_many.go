package fileutil

import (
	"os"

	"github.com/blend/go-sdk/ex"
)

// RemoveMany removes an array of files.
func RemoveMany(filePaths ...string) error {
	var err error
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			return ex.New(err, ex.OptMessage(path))
		}
	}
	return err
}
