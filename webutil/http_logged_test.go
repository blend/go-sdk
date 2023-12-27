/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestHTTPLogged(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := logger.New(
		logger.OptOutput(buf),
		logger.OptAll(),
	)
	assert.Nil(err)

	var didCall bool
	server := httptest.NewServer(NestMiddleware(func(rw http.ResponseWriter, req *http.Request) {
		didCall = true
	}, HTTPLogged(log)))

	res, err := http.Get(server.URL)
	assert.Nil(err)
	defer res.Body.Close()
	assert.True(didCall)
}

func TestFormatHeaders(t *testing.T) {
	assert := assert.New(t)

	tf := logger.NewTextOutputFormatter(logger.OptTextNoColor())
	actual := FormatHeaders(tf, ansi.ColorBlue, http.Header{"Foo": []string{"bar"}, "Moo": []string{"loo"}})
	assert.Equal("{ Foo:bar Moo:loo }", actual)

	actual = FormatHeaders(tf, ansi.ColorBlue, http.Header{"Moo": []string{"loo"}, "Foo": []string{"bar"}})
	assert.Equal("{ Foo:bar Moo:loo }", actual)

	tf = logger.NewTextOutputFormatter()
	actual = FormatHeaders(tf, ansi.ColorBlue, http.Header{"Foo": []string{"bar"}, "Moo": []string{"loo"}})
	assert.Equal("{ "+ansi.ColorBlue.Apply("Foo")+":bar "+ansi.ColorBlue.Apply("Moo")+":loo }", actual)
}

func TestColorizeByStatusCode(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		StatusCode	int
		Value		string
		Expected	string
	}{
		{StatusCode: http.StatusInternalServerError, Value: "this is a server error", Expected: ansi.ColorRed.Apply("this is a server error")},
		{StatusCode: http.StatusBadRequest, Value: "this is a bad request", Expected: ansi.ColorYellow.Apply("this is a bad request")},
		{StatusCode: http.StatusOK, Value: "this is ok", Expected: ansi.ColorGreen.Apply("this is ok")},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, ColorizeByStatusCode(tc.StatusCode, tc.Value))
	}
}

func TestColorizebyStatusCodeWithFormatter(t *testing.T) {
	assert := assert.New(t)

	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}
	color := logger.TextOutputFormatter{
		NoColor: false,
	}

	testCases := [...]struct {
		StatusCode	int
		Formatter	logger.TextFormatter
		Value		string
		Expected	string
	}{
		// Color
		{StatusCode: http.StatusInternalServerError, Value: "this is a server error", Formatter: color, Expected: ansi.ColorRed.Apply("this is a server error")},
		{StatusCode: http.StatusBadRequest, Value: "this is a bad request", Formatter: color, Expected: ansi.ColorYellow.Apply("this is a bad request")},
		{StatusCode: http.StatusOK, Value: "this is ok", Formatter: color, Expected: ansi.ColorGreen.Apply("this is ok")},

		// NoColor
		{StatusCode: http.StatusInternalServerError, Value: "this is a server error", Formatter: noColor, Expected: "this is a server error"},
		{StatusCode: http.StatusBadRequest, Value: "this is a bad request", Formatter: noColor, Expected: "this is a bad request"},
		{StatusCode: http.StatusOK, Value: "this is ok", Formatter: noColor, Expected: "this is ok"},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, ColorizeByStatusCodeWithFormatter(tc.Formatter, tc.StatusCode, tc.Value))
	}
}

func TestColorizeStatusCode(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		StatusCode	int
		Expected	string
	}{
		{StatusCode: http.StatusInternalServerError, Expected: ansi.ColorRed.Apply("500")},
		{StatusCode: http.StatusBadRequest, Expected: ansi.ColorYellow.Apply("400")},
		{StatusCode: http.StatusOK, Expected: ansi.ColorGreen.Apply("200")},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, ColorizeStatusCode(tc.StatusCode))
	}
}

func TestColorizeStatusCodeWithFormatter(t *testing.T) {
	assert := assert.New(t)

	noColor := logger.TextOutputFormatter{
		NoColor: true,
	}
	color := logger.TextOutputFormatter{
		NoColor: false,
	}

	testCases := [...]struct {
		StatusCode	int
		Formatter	logger.TextFormatter
		Expected	string
	}{
		// Color
		{StatusCode: http.StatusInternalServerError, Formatter: color, Expected: ansi.ColorRed.Apply("500")},
		{StatusCode: http.StatusBadRequest, Formatter: color, Expected: ansi.ColorYellow.Apply("400")},
		{StatusCode: http.StatusOK, Formatter: color, Expected: ansi.ColorGreen.Apply("200")},

		// NoColor
		{StatusCode: http.StatusInternalServerError, Formatter: noColor, Expected: "500"},
		{StatusCode: http.StatusBadRequest, Formatter: noColor, Expected: "400"},
		{StatusCode: http.StatusOK, Formatter: noColor, Expected: "200"},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, ColorizeStatusCodeWithFormatter(tc.Formatter, tc.StatusCode))
	}
}
