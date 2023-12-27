/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"io"
	"net/http"

	"github.com/blend/go-sdk/ex"
)

const (
	// DefaultPostedFilesMaxMemory is the maximum post body size we will typically consume.
	DefaultPostedFilesMaxMemory = 67_108_864	//64mb
)

// PostedFile is a file that has been posted to an hc endpoint.
type PostedFile struct {
	Key		string
	FileName	string
	Contents	[]byte
	ContentType	string
}

// PostedFilesOptions are options for the PostedFiles function.
type PostedFilesOptions struct {
	MaxMemory		int64
	ParseMultipartForm	bool
	ParseForm		bool
}

// PostedFileOption mutates posted file options.
type PostedFileOption func(*PostedFilesOptions)

// OptPostedFilesMaxMemory sets the max memory for the posted files options (defaults to 64mb).
func OptPostedFilesMaxMemory(maxMemory int64) PostedFileOption {
	return func(pfo *PostedFilesOptions) { pfo.MaxMemory = maxMemory }
}

// OptPostedFilesParseMultipartForm sets if we should parse the multipart form for files (defaults to true).
func OptPostedFilesParseMultipartForm(parseMultipartForm bool) PostedFileOption {
	return func(pfo *PostedFilesOptions) { pfo.ParseMultipartForm = parseMultipartForm }
}

// OptPostedFilesParseForm sets if we should parse the post form for files (defaults to false).
func OptPostedFilesParseForm(parseForm bool) PostedFileOption {
	return func(pfo *PostedFilesOptions) { pfo.ParseForm = parseForm }
}

// PostedFiles returns any files posted to the request.
//
// The files are held in memory, if you need to stream the files out because they may be large,
// you should use the `*net/http.Request.FormFile(...)` function directly instead of this method.
func PostedFiles(r *http.Request, opts ...PostedFileOption) ([]PostedFile, error) {
	var files []PostedFile

	options := PostedFilesOptions{
		MaxMemory:		DefaultPostedFilesMaxMemory,
		ParseMultipartForm:	true,
		ParseForm:		false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.ParseMultipartForm {
		if err := r.ParseMultipartForm(options.MaxMemory); err != nil {
			return nil, ex.New(err)
		}
		for key := range r.MultipartForm.File {
			file, err := readPostedFile(r, key)
			if err != nil {
				return nil, err
			}
			files = append(files, *file)
		}
	}
	if options.ParseForm {
		if err := r.ParseForm(); err != nil {
			return nil, ex.New(err)
		}
		for key := range r.PostForm {
			file, err := readPostedFile(r, key)
			if err != nil {
				return nil, err
			}
			files = append(files, *file)
		}
	}
	return files, nil
}

func readPostedFile(r *http.Request, key string) (*PostedFile, error) {
	fileReader, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, ex.New(err)
	}
	defer fileReader.Close()
	fileContents, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, ex.New(err)
	}
	return &PostedFile{Key: key, FileName: fileHeader.Filename, Contents: fileContents}, nil
}
