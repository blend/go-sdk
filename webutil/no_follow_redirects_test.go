package webutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNoFollowRedirects(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(http.ErrUseLastResponse, NoFollowRedirects()(nil, nil))

	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!")
		return
	}))
	defer second.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, second.URL, 307)
		return
	}))
	defer server.Close()

	client := http.Client{
		CheckRedirect: NoFollowRedirects(),
	}

	res, err := client.Get(server.URL)
	assert.Nil(err)
	defer res.Body.Close()
	assert.Equal(307, res.StatusCode, "the redirect status code should be returned by the server")
}
