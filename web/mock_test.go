/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

func TestMock(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.GET("/", func(_ *Ctx) Result { return NoContent })

	res, err := Mock(app, &http.Request{Method: "GET", URL: &url.URL{Scheme: webutil.SchemeHTTP, Path: "/"}}).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)

	assert.True(app.IsStopped())

	// try to make another request to the underlying test server

	res, err = http.Get(res.Request.URL.String())
	assert.NotNil(err)
	assert.Nil(res)
}

func TestMockGet(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.GET("/", func(_ *Ctx) Result { return NoContent })

	res, err := MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)

	assert.True(app.IsStopped())
}

func TestMockPostedFile(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.POST("/", func(r *Ctx) Result {
		postedFiles, err := webutil.PostedFiles(r.Request)
		if err != nil {
			return Text.BadRequest(err)
		}
		if len(postedFiles) == 0 {
			return Text.BadRequest(fmt.Errorf("there should be 2 files"))
		}
		return Text.OK()
	})

	res, err := MockMethod(app, http.MethodPost, "/",
		r2.OptPostedFiles(
			webutil.PostedFile{Key: "file0", FileName: "file0.txt", Contents: []byte("this is just a test")},
			webutil.PostedFile{Key: "file1", FileName: "file1.txt", Contents: []byte("this is just a test")},
		),
	).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}
