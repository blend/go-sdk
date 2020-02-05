package webutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNoFollowRedirects(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(http.ErrUseLastResponse, NoFollowRedirects()(nil, nil))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bad", 307)
		return
	}))
	defer server.Close()

	client := http.Client{
		CheckRedirect: NoFollowRedirects(),
	}

	res, err := client.Get(server.URL)
	assert.Nil(err)
	defer res.Body.Close()
	assert.Equal(307, res.StatusCode)
}
