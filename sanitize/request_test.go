package sanitize

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSanitizeRequest(t *testing.T) {
	it := assert.New(t)

	req := &http.Request{
		Header: http.Header{
			"Accept":         {"application/json"},
			"Authorization":  {"Bearer foo", "Bearer bar"},
			"X-Secret-Token": {"super_secret_token"},
		},
		URL: &url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/api/sensitive",
			RawQuery: (url.Values{
				"ok":           {"ok0", "ok1"},
				"access_token": {"super_secret"},
				"sensitive":    {"sensitive0", "sensitive1"},
			}).Encode(),
		},
	}

	output := Request(req,
		OptRequestAddDisallowedHeaders("X-Secret-Token"),
		OptRequestAddDisallowedQueryParams("sensitive"),
		OptRequestValueSanitizer(func(key string, values ...string) []string {
			return []string{"***"}
		}),
	)

	it.NotNil(output)
	it.Equal([]string{"application/json"}, req.Header["Accept"])
	it.Equal([]string{"***"}, output.Header["Authorization"])
	it.Equal([]string{"Bearer foo", "Bearer bar"}, req.Header["Authorization"])
	it.Equal([]string{"***"}, output.Header["X-Secret-Token"])

	it.Equal([]string{"***"}, output.URL.Query()["access_token"])
	it.Equal([]string{"***"}, output.URL.Query()["sensitive"])
}
