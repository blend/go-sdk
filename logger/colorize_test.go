package logger

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestFlagTextColor(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Flag     string
		Expected ansi.Color
	}{
		{Info, ansi.ColorLightWhite},
		{Debug, ansi.ColorLightYellow},
		{Warning, ansi.ColorLightYellow},
		{Error, ansi.ColorRed},
		{Fatal, ansi.ColorRed},
		{"foo", DefaultFlagTextColor},
		{"", DefaultFlagTextColor},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, FlagTextColor(tc.Flag))
	}
}

func TestColorizeByStatusCode(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		StatusCode int
		Value      string
		Expected   string
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

	noColor := TextOutputFormatter{
		NoColor: true,
	}
	color := TextOutputFormatter{
		NoColor: false,
	}

	testCases := [...]struct {
		StatusCode int
		Formatter  TextFormatter
		Value      string
		Expected   string
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
		StatusCode int
		Expected   string
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

	noColor := TextOutputFormatter{
		NoColor: true,
	}
	color := TextOutputFormatter{
		NoColor: false,
	}

	testCases := [...]struct {
		StatusCode int
		Formatter  TextFormatter
		Expected   string
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
