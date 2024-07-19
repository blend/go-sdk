/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// Webhook is a configurable request.
type Webhook struct {
	Method  string            `json:"method" yaml:"method"`
	URL     string            `json:"url" yaml:"url"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    string            `json:"body" yaml:"body"`
}

// IsZero returns if the webhook is set.
func (wh Webhook) IsZero() bool {
	return wh.URL == ""
}

// MethodOrDefault returns the method or a default.
func (wh Webhook) MethodOrDefault() string {
	if wh.Method != "" {
		return wh.Method
	}
	return "GET"
}

// Send sends the webhook.
func (wh Webhook) Send() (*http.Response, error) {
	u, err := url.Parse(wh.URL)
	if err != nil {
		return nil, ex.New(err)
	}

	req := &http.Request{
		Method: wh.MethodOrDefault(),
		URL:    u,
	}
	headers := http.Header{}
	for key, value := range wh.Headers {
		headers.Add(key, value)
	}
	req.Header = headers
	if wh.Body != "" {
		req.Body = io.NopCloser(bytes.NewBufferString(wh.Body))
		req.ContentLength = int64(len(wh.Body))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ex.New(err)
	}
	return res, nil
}
