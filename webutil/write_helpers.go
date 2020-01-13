package webutil

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"github.com/blend/go-sdk/ex"
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
	return ex.New(json.NewEncoder(w).Encode(response))
}

// WriteXML marshalls an object to json.
func WriteXML(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set(HeaderContentType, ContentTypeXML)
	w.WriteHeader(statusCode)
	return ex.New(xml.NewEncoder(w).Encode(response))
}

// DeserializeReaderAsJSON deserializes a post body as json to a given object.
func DeserializeReaderAsJSON(object interface{}, body io.ReadCloser) error {
	defer body.Close()
	return ex.New(json.NewDecoder(body).Decode(object))
}
