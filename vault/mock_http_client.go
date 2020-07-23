package vault

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// NewMockHTTPClient returns a new mock http client.
// MockHTTPClient is used to test VaultClient itself, and should
// not be used for your own mocks.
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		contents: make(map[string]*http.Response),
	}
}

// MockHTTPClient is a mock http client.
// It is used to test the vault client iself, and should not be used for your own mocks.
type MockHTTPClient struct {
	contents map[string]*http.Response
}

// With adds a mocked endpoint.
func (mh *MockHTTPClient) With(verb string, url *url.URL, response *http.Response) *MockHTTPClient {
	mh.contents[fmt.Sprintf("%s_%s", verb, url.String())] = response
	return mh
}

// WithString adds a mocked endpoint.
func (mh *MockHTTPClient) WithString(verb string, url *url.URL, contents string) *MockHTTPClient {
	mh.contents[fmt.Sprintf("%s_%s", verb, url.String())] = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(contents))),
	}
	return mh
}

// Do implements HTTPClient.
func (mh *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if res, hasRes := mh.contents[fmt.Sprintf("%s_%s", req.Method, req.URL.String())]; hasRes {
		return res, nil
	}
	return nil, fmt.Errorf("not found")
}
