package r2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptJSONBody(t *testing.T) {
	assert := assert.New(t)

	object := map[string]interface{}{"foo": "bar"}

	opt := OptJSONBody(object)

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Body)

	contents, err := ioutil.ReadAll(req.Request.Body)
	assert.Nil(err)
	assert.Equal(`{"foo":"bar"}`, string(contents))

	assert.Equal(ContentTypeApplicationJSON, req.Request.Header.Get("Content-Type"))
}

func TestOptJSONBodyRedirect(t *testing.T) {
	assert := assert.New(t)

	finalServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 13 {
			http.Error(rw, fmt.Sprintf("final; invalid content length: %d", r.ContentLength), http.StatusBadRequest)
			return
		}
		var actualBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&actualBody); err != nil {
			http.Error(rw, "final; invalid json body", http.StatusBadRequest)
			return
		}
		if actualBody["foo"] != "bar" {
			http.Error(rw, "final; invalid foo", http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "OK!")
	}))
	defer finalServer.Close()

	// we test the redirect to assert that the .GetBody function works as intended
	redirectServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 13 {
			http.Error(rw, fmt.Sprintf("redirect; invalid content length: %d", r.ContentLength), http.StatusBadRequest)
			return
		}
		// http.StatusTemporaryRedirect necessary here so the body follows the redirect
		http.Redirect(rw, r, finalServer.URL, http.StatusTemporaryRedirect)
	}))
	defer redirectServer.Close()

	jsonObject := map[string]interface{}{"foo": "bar"}

	r := New(redirectServer.URL,
		OptJSONBody(jsonObject),
	)
	assert.Equal(13, r.Request.ContentLength)
	assert.NotNil(r.Request.Body)
	assert.NotNil(r.Request.GetBody)

	contents, meta, err := r.Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
}
