/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

type xmlBody struct {
	X []string `xml:"x"`
	Y []string `xml:"y"`
}

func TestRequestOptions(t *testing.T) {
	assert := assert.New(t)

	req := &http.Request{}

	assert.Empty(req.Method)
	assert.Nil(OptMethod("POST")(req))
	assert.Equal("POST", req.Method)

	req.Method = ""
	assert.Nil(OptGet()(req))
	assert.Equal("GET", req.Method)

	req.Method = ""
	assert.Nil(OptPost()(req))
	assert.Equal("POST", req.Method)

	req.Method = ""
	assert.Nil(OptPut()(req))
	assert.Equal("PUT", req.Method)

	req.Method = ""
	assert.Nil(OptPatch()(req))
	assert.Equal("PATCH", req.Method)

	req.Method = ""
	assert.Nil(OptDelete()(req))
	assert.Equal("DELETE", req.Method)

	type contextKey struct{}
	assert.Nil(req.Context().Value(contextKey{}))
	assert.Nil(OptContext(context.WithValue(context.Background(), contextKey{}, "foo"))(req))
	assert.Equal("foo", req.Context().Value(contextKey{}))

	assert.Nil(req.URL)
	assert.Nil(OptQuery(url.Values{"foo": []string{"bar", "baz"}})(req))
	assert.NotNil(req.URL)
	assert.Equal("foo=bar&foo=baz", req.URL.RawQuery)

	req.URL = &url.URL{}
	assert.Nil(OptQueryValue("foo", "bar")(req))
	assert.NotNil(req.URL)
	assert.Equal("foo=bar", req.URL.RawQuery)

	assert.Nil(req.Header)
	assert.Nil(OptHeader(http.Header{"X-Foo": []string{"bar", "baz"}})(req))
	assert.Equal("bar", req.Header.Get("X-Foo"))

	req.Header = nil
	assert.Nil(OptHeaderValue("X-Foo", "bar")(req))
	assert.Equal("bar", req.Header.Get("X-Foo"))

	assert.Nil(req.PostForm)
	assert.Nil(OptPostForm(url.Values{"foo": []string{"bar", "baz"}})(req))
	assert.Equal("bar", req.PostForm.Get("foo"))

	req.PostForm = nil
	assert.Nil(OptPostFormValue("buzz", "fuzz")(req))
	assert.Equal("fuzz", req.PostForm.Get("buzz"))

	req.Header = nil
	assert.Nil(OptCookie(&http.Cookie{Name: "sid", Value: "my value"})(req))
	c, err := req.Cookie("sid")
	assert.Nil(err)
	assert.Equal("my value", c.Value)

	req.Header = nil
	assert.Nil(OptCookieValue("jsid", "another value")(req))
	c, err = req.Cookie("jsid")
	assert.Nil(err)
	assert.Equal("another value", c.Value)

	assert.Nil(req.Body)
	assert.Nil(OptBody(io.NopCloser(bytes.NewReader([]byte("foo bar"))))(req))
	assert.NotNil(req.Body)
	read, err := io.ReadAll(req.Body)
	assert.Nil(err)
	assert.Equal([]byte("foo bar"), read)

	req.Body = nil
	assert.Nil(OptBodyBytes([]byte("bar foo"))(req))
	assert.NotNil(req.Body)
	read, err = io.ReadAll(req.Body)
	assert.Nil(err)
	assert.Equal([]byte("bar foo"), read)

	postedFiles := []PostedFile{
		{Key: "file0", FileName: "file.txt", Contents: []byte("foo bar baz")},
		{Key: "file1", FileName: "file_1.txt", Contents: []byte("fuzzy wuzzy was a bear")},
	}
	req.Header = nil
	req.Body = nil
	assert.Nil(OptPostedFiles(postedFiles...)(req))
	assert.NotEmpty(req.Header)
	assert.NotNil(req.Body)

	req.Header = nil
	req.Body = nil
	assert.Nil(OptJSONBody([]string{"foo", "bar"})(req))
	assert.Equal(ContentTypeApplicationJSON, req.Header.Get(HeaderContentType))
	assert.NotNil(req.Body)

	req.Header = nil
	req.Body = nil
	assert.Nil(OptXMLBody([]string{"foo", "bar"})(req))
	assert.Equal(ContentTypeApplicationXML, req.Header.Get(HeaderContentType))
	assert.NotNil(req.Body)
}

func TestOptBodyBytes(t *testing.T) {
	assert := assert.New(t)
	body := []byte("hello\n")
	opt := OptBodyBytes(body)

	r := &http.Request{}
	err := opt(r)
	assert.Nil(err)

	bodyBytes, err := io.ReadAll(r.Body)
	assert.Nil(err)
	assert.Equal(body, bodyBytes)
	assert.Equal(r.ContentLength, 6)
	validateRequestGetBody(assert, r, body)
}

func TestOptPostedFiles(t *testing.T) {
	assert := assert.New(t)
	file1 := PostedFile{Key: "a", FileName: "b.txt", Contents: []byte("hey")}
	file2 := PostedFile{Key: "c", FileName: "d.txt", Contents: []byte("bye")}
	file3 := PostedFile{Key: "d", FileName: "e.txt", Contents: []byte("bytes"), ContentType: "application/pdf"}
	opt := OptPostedFiles(file1, file2, file3)

	r := &http.Request{}
	err := opt(r)
	assert.Nil(err)

	boundary := getBoundary(assert, r.Header)
	ct := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)
	assert.Equal(r.Header, http.Header{HeaderContentType: []string{ct}})
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	assert.Nil(err)
	expected := fmt.Sprintf(
		"--%[1]s\r\nContent-Disposition: form-data; name=%[2]q; filename=%[3]q\r\n"+
			"Content-Type: application/octet-stream\r\n\r\n%[4]s\r\n"+
			"--%[1]s\r\nContent-Disposition: form-data; name=%[5]q; filename=%[6]q\r\n"+
			"Content-Type: application/octet-stream\r\n\r\n%[7]s\r\n"+
			"--%[1]s\r\nContent-Disposition: form-data; name=%[8]q; filename=%[9]q\r\n"+
			"Content-Type: %[10]s\r\n\r\n%[11]s\r\n--%[1]s--\r\n",
		boundary,
		file1.Key,
		file1.FileName,
		file1.Contents,
		file2.Key,
		file2.FileName,
		file2.Contents,
		file3.Key,
		file3.FileName,
		file3.ContentType,
		file3.Contents,
	)
	assert.Equal(expected, string(bodyBytes))
	assert.Equal(r.ContentLength, len(expected))
	validateRequestGetBody(assert, r, []byte(expected))
}

func TestOptPostedFiles_postForm(t *testing.T) {
	assert := assert.New(t)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.FormValue("foo") != "foo-value" {
			http.Error(w, "post form value `foo` missing or incorrect", http.StatusBadRequest)
			return
		}
		if req.FormValue("bar") != "bar-value" {
			http.Error(w, "post form value `bar` missing or incorrect", http.StatusBadRequest)
			return
		}

		files, err := PostedFiles(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(files) < 2 {
			http.Error(w, "posted files missing", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!")
	}))
	defer testServer.Close()

	file1 := PostedFile{Key: "a", FileName: "b.txt", Contents: []byte("hey")}
	file2 := PostedFile{Key: "c", FileName: "d.txt", Contents: []byte("bye")}

	postFormOpt0 := OptPostFormValue("foo", "foo-value")
	postFormOpt1 := OptPostFormValue("bar", "bar-value")
	opt := OptPostedFiles(file1, file2)

	r, err := http.NewRequest(http.MethodPost, testServer.URL, nil)
	assert.Nil(err)
	err = postFormOpt0(r)
	assert.Nil(err)
	err = postFormOpt1(r)
	assert.Nil(err)
	err = opt(r)
	assert.Nil(err)

	res, err := http.DefaultClient.Do(r)
	assert.Nil(err)
	defer res.Body.Close()
	message, _ := io.ReadAll(res.Body)
	assert.Equal(http.StatusOK, res.StatusCode, string(message))
}

func TestOptJSONBody(t *testing.T) {
	assert := assert.New(t)
	payload := map[string]float64{"x": 1.25, "y": -5.75}
	opt := OptJSONBody(payload)

	r := &http.Request{}
	err := opt(r)
	assert.Nil(err)

	assert.NotNil(r.Body)
	assert.NotNil(r.GetBody)

	assert.Equal(r.Header, http.Header{HeaderContentType: []string{ContentTypeApplicationJSON}})
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	assert.Nil(err)
	expected := []byte(`{"x":1.25,"y":-5.75}`)
	assert.Equal(expected, bodyBytes)
	assert.Equal(r.ContentLength, 20)

	validateRequestGetBody(assert, r, expected)

	verifyContents, _ := json.Marshal(payload)
	verify, err := http.NewRequest("POST", "https://example.com/test", bytes.NewReader(verifyContents))
	assert.Nil(err)
	assert.Equal(verify.ContentLength, r.ContentLength)
}

func TestOptXMLBody(t *testing.T) {
	assert := assert.New(t)
	payload := xmlBody{X: []string{"hello"}, Y: []string{"goodbye"}}
	opt := OptXMLBody(payload)

	r := &http.Request{}
	err := opt(r)
	assert.Nil(err)

	assert.Equal(r.Header, http.Header{HeaderContentType: []string{ContentTypeApplicationXML}})
	bodyBytes, err := io.ReadAll(r.Body)
	assert.Nil(err)
	expected := []byte("<xmlBody><x>hello</x><y>goodbye</y></xmlBody>")
	assert.Equal(expected, bodyBytes)
	assert.Equal(r.ContentLength, 45)
	validateRequestGetBody(assert, r, expected)
}

func getBoundary(assert *assert.Assertions, h http.Header) string {
	boundaryPrefix := "multipart/form-data; boundary="
	ct := h.Get(HeaderContentType)
	assert.True(strings.HasPrefix(ct, boundaryPrefix))
	return strings.TrimPrefix(ct, boundaryPrefix)
}

func validateRequestGetBody(assert *assert.Assertions, r *http.Request, expected []byte) {
	assert.NotNil(r.GetBody)
	bodyRC, err := r.GetBody()
	assert.Nil(err)
	defer bodyRC.Close()
	bodyBytes, err := io.ReadAll(bodyRC)
	assert.Nil(err)
	assert.Equal(expected, bodyBytes)
}
