package web

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// BaseHeaders are the default headers added by go-web.
func BaseHeaders() http.Header {
	return http.Header{
		webutil.HeaderServer: []string{PackageName},
	}
}
