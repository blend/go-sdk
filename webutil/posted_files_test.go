/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_PostedFiles(t *testing.T) {
	its := assert.New(t)

	file0 := PostedFile{
		Key:		"file0",
		FileName:	"file0.txt",
		Contents:	[]byte("file0-contents"),
	}
	file1 := PostedFile{
		Key:		"file1",
		FileName:	"file1.txt",
		Contents:	[]byte(strings.Repeat("a", 1<<20)),
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		files, err := PostedFiles(req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if len(files) != 2 {
			http.Error(rw, "invalid file count", http.StatusBadRequest)
			return
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Key < files[j].Key
		})

		if files[0].Key != file0.Key {
			http.Error(rw, fmt.Sprintf("invalid file0 key: %s, expectd: %s", files[0].Key, file0.Key), http.StatusBadRequest)
			return
		}
		if files[0].FileName != file0.FileName {
			http.Error(rw, fmt.Sprintf("invalid file0 key: %s, expectd: %s", files[0].FileName, file0.FileName), http.StatusBadRequest)
			return
		}
		if string(files[0].Contents) != string(file0.Contents) {
			http.Error(rw, fmt.Sprintf("invalid file0 contents: %s, expectd: %s", files[0].Contents, file0.Contents), http.StatusBadRequest)
			return
		}

		if files[1].Key != file1.Key {
			http.Error(rw, fmt.Sprintf("invalid file1 key: %s, expectd: %s", files[1].Key, file1.Key), http.StatusBadRequest)
			return
		}
		if files[1].FileName != file1.FileName {
			http.Error(rw, fmt.Sprintf("invalid file1 key: %s, expectd: %s", files[1].FileName, file1.FileName), http.StatusBadRequest)
			return
		}
		if string(files[1].Contents) != string(file1.Contents) {
			http.Error(rw, "invalid file1 contents", http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!")
		return
	}))
	defer server.Close()

	r, err := http.NewRequest(http.MethodPost, server.URL, nil)
	its.Nil(err)
	err = OptPostedFiles(file0, file1)(r)
	its.Nil(err)
	res, err := http.DefaultClient.Do(r)
	its.Nil(err)
	defer res.Body.Close()
	bodyContents, _ := io.ReadAll(res.Body)
	its.Equal(http.StatusOK, res.StatusCode, string(bodyContents))
}

func Test_PostedFiles_onlyParseForm(t *testing.T) {
	its := assert.New(t)

	file0 := PostedFile{
		Key:		"file0",
		FileName:	"file0.txt",
		Contents:	[]byte("file0-contents"),
	}
	file1 := PostedFile{
		Key:		"file1",
		FileName:	"file1.txt",
		Contents:	[]byte(strings.Repeat("a", 1<<20)),
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		files, err := PostedFiles(req,
			OptPostedFilesParseMultipartForm(false),
			OptPostedFilesParseForm(true),
		)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if len(files) != 0 {
			http.Error(rw, fmt.Sprintf("invalid file count: %d", len(files)), http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!")
		return
	}))
	defer server.Close()

	r, err := http.NewRequest(http.MethodPost, server.URL, nil)
	its.Nil(err)
	err = OptPostedFiles(file0, file1)(r)
	its.Nil(err)
	res, err := http.DefaultClient.Do(r)
	its.Nil(err)
	defer res.Body.Close()
	contents, err := io.ReadAll(res.Body)
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode, string(contents))
}

func Test_PostedFiles_maxMemory(t *testing.T) {
	its := assert.New(t)

	file0 := PostedFile{
		Key:		"file0",
		FileName:	"file0.txt",
		Contents:	[]byte("file0-contents"),
	}
	file1 := PostedFile{
		Key:		"file1",
		FileName:	"file1.txt",
		Contents:	[]byte(strings.Repeat("a", 1<<20)),
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		files, err := PostedFiles(req,
			OptPostedFilesMaxMemory(1<<10),
			OptPostedFilesParseMultipartForm(true),
			OptPostedFilesParseForm(false),
		)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if len(files) != 2 {
			http.Error(rw, fmt.Sprintf("invalid file count: %d", len(files)), http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!")
		return
	}))
	defer server.Close()

	r, err := http.NewRequest(http.MethodPost, server.URL, nil)
	its.Nil(err)
	err = OptPostedFiles(file0, file1)(r)
	its.Nil(err)
	res, err := http.DefaultClient.Do(r)
	its.Nil(err)
	defer res.Body.Close()
	contents, err := io.ReadAll(res.Body)
	its.Nil(err)
	its.Equal(http.StatusOK, res.StatusCode, string(contents))
}
