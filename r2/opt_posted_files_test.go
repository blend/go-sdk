/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestOptPostedFiles(t *testing.T) {
	its := assert.New(t)

	r := New(TestURL, OptPostedFiles(
		webutil.PostedFile{
			Key:		"form-key",
			FileName:	"file.txt",
			Contents:	[]byte("this is a test"),
		},
		webutil.PostedFile{
			Key:		"form-key-2",
			FileName:	"file2.txt",
			Contents:	[]byte("this is a test2"),
		},
	))
	its.NotNil(r.Request.Body)

	files, err := webutil.PostedFiles(r.Request)
	its.Nil(err)
	its.Len(files, 2)

	its.AnyCount(files, 1, func(v interface{}) bool {
		file := v.(webutil.PostedFile)
		return file.Key == "form-key" &&
			file.FileName == "file.txt" &&
			string(file.Contents) == "this is a test"
	})

	its.AnyCount(files, 1, func(v interface{}) bool {
		file := v.(webutil.PostedFile)
		return file.Key == "form-key-2" &&
			file.FileName == "file2.txt" &&
			string(file.Contents) == "this is a test2"
	})
}
