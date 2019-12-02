package webutil

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
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
	assert.Nil(OptBody(ioutil.NopCloser(bytes.NewReader([]byte("foo bar"))))(req))
	assert.NotNil(req.Body)
	read, err := ioutil.ReadAll(req.Body)
	assert.Nil(err)
	assert.Equal([]byte("foo bar"), read)

	req.Body = nil
	assert.Nil(OptBodyBytes([]byte("bar foo"))(req))
	assert.NotNil(req.Body)
	read, err = ioutil.ReadAll(req.Body)
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

	bodyBytes, err := ioutil.ReadAll(r.Body)
	assert.Nil(err)
	assert.Equal(bodyBytes, body)
	assert.Equal(r.ContentLength, 6)
}

func TestOptJSONBody(t *testing.T) {
	assert := assert.New(t)
	payload := map[string]float64{"x": 1.25, "y": -5.75}
	opt := OptJSONBody(payload)

	r := &http.Request{}
	err := opt(r)
	assert.Nil(err)

	assert.Equal(r.Header, http.Header{HeaderContentType: []string{ContentTypeApplicationJSON}})
	bodyBytes, err := ioutil.ReadAll(r.Body)
	assert.Nil(err)
	assert.Equal(bodyBytes, []byte(`{"x":1.25,"y":-5.75}`))
	assert.Equal(r.ContentLength, 20)
}

func TestOptXMLBody(t *testing.T) {
	assert := assert.New(t)
	payload := xmlBody{X: []string{"hello"}, Y: []string{"goodbye"}}
	opt := OptXMLBody(payload)

	r := &http.Request{}
	err := opt(r)
	assert.Nil(err)

	assert.Equal(r.Header, http.Header{HeaderContentType: []string{ContentTypeApplicationXML}})
	bodyBytes, err := ioutil.ReadAll(r.Body)
	assert.Nil(err)
	assert.Equal(bodyBytes, []byte("<xmlBody><x>hello</x><y>goodbye</y></xmlBody>"))
	assert.Equal(r.ContentLength, 45)
}
