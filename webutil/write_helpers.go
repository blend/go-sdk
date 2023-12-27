/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net"
	"net/http"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrNetWrite ex.Class = "network write error"
)

// WriteNoContent writes http.StatusNoContent for a request.
func WriteNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// WriteRawContent writes raw content for the request.
func WriteRawContent(w http.ResponseWriter, statusCode int, content []byte) error {
	w.WriteHeader(statusCode)
	_, err := w.Write(content)
	return ex.New(err)
}

// WriteJSON marshalls an object to json.
func WriteJSON(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		if typed, ok := err.(*net.OpError); ok {
			return ex.New(ErrNetWrite, ex.OptInner(typed))
		}
		return ex.New(err)
	}
	return nil
}

// WriteXML marshalls an object to json.
func WriteXML(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set(HeaderContentType, ContentTypeXML)
	w.WriteHeader(statusCode)
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		if typed, ok := err.(*net.OpError); ok {
			return ex.New(ErrNetWrite, ex.OptInner(typed))
		}
		return ex.New(err)
	}
	return nil
}

// DeserializeReaderAsJSON deserializes a post body as json to a given object.
func DeserializeReaderAsJSON(object interface{}, body io.ReadCloser) error {
	defer body.Close()
	return ex.New(json.NewDecoder(body).Decode(object))
}
