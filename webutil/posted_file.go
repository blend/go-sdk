package webutil

import (
	"io/ioutil"
	"net/http"
	"github.com/blend/go-sdk/ex"
)

const (
	// MaxPostBodySize is the maximum post body size we will typically consume.
	MaxPostBodySize = int64(1 << 26) //64mb
)

// PostedFile is a file that has been posted to an hc endpoint.
type PostedFile struct {
	Key      string
	FileName string
	Contents []byte
}

// PostedFiles returns any files posted
func PostedFiles(r *http.Request) ([]PostedFile, error) {
	var files []PostedFile

	err := r.ParseMultipartForm(MaxPostBodySize)
	if err == nil {
		for key := range r.MultipartForm.File {
			fileReader, fileHeader, err := r.FormFile(key)
			if err != nil {
				return nil, ex.New(err)
			}
			bytes, err := ioutil.ReadAll(fileReader)
			if err != nil {
				return nil, ex.New(err)
			}
			files = append(files, PostedFile{Key: key, FileName: fileHeader.Filename, Contents: bytes})
		}
	} else {
		err = r.ParseForm()
		if err == nil {
			for key := range r.PostForm {
				if fileReader, fileHeader, err := r.FormFile(key); err == nil && fileReader != nil {
					bytes, err := ioutil.ReadAll(fileReader)
					if err != nil {
						return nil, ex.New(err)
					}
					files = append(files, PostedFile{Key: key, FileName: fileHeader.Filename, Contents: bytes})
				}
			}
		}
	}
	return files, nil
}
