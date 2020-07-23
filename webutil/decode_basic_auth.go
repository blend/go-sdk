package webutil

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// DecodeBasicAuth pulls a basic auth header from a request and returns the
// username and password that were passed.
func DecodeBasicAuth(req *http.Request) (username, password string, err error) {
	var rawHeader string
	rawHeader, err = headerValue(HeaderAuthorization, req)
	if err != nil {
		err = ex.New(ErrUnauthorized, ex.OptInner(err))
		return
	}

	auth := strings.SplitN(rawHeader, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		err = ex.New(ErrUnauthorized)
		return
	}

	var payload []byte
	payload, err = base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		err = ex.New(ErrUnauthorized, ex.OptInner(err))
		return
	}

	username, password, err = SplitColon(string(payload))
	if err != nil {
		err = ex.New(ErrUnauthorized, ex.OptInner(err))
		return
	}

	return
}

func headerValue(key string, req *http.Request) (value string, err error) {
	if value = req.Header.Get(key); len(value) > 0 {
		return
	}
	err = ErrParameterMissing
	return
}
