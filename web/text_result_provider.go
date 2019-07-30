package web

import (
	"fmt"
	"net/http"
)

var (
	// Text is a static singleton text result provider.
	Text TextResultProvider

	// assert TestResultProvider implements result provider.
	_ ResultProvider = Text
)

// TextResultProvider is the default response provider if none is specified.
type TextResultProvider struct{}

// NotFound returns a plaintext result.
func (trp TextResultProvider) NotFound() Result {
	return &RawResult{
		StatusCode:  http.StatusNotFound,
		ContentType: ContentTypeText,
		Response:    []byte("Not Found"),
	}
}

// NotAuthorized returns a plaintext result.
func (trp TextResultProvider) NotAuthorized() Result {
	return &RawResult{
		StatusCode:  http.StatusUnauthorized,
		ContentType: ContentTypeText,
		Response:    []byte("Not Authorized"),
	}
}

// InternalError returns a plainttext result.
func (trp TextResultProvider) InternalError(err error) Result {
	return ResultWithLoggedError(&RawResult{
		StatusCode:  http.StatusInternalServerError,
		ContentType: ContentTypeText,
		Response:    []byte(fmt.Sprintf("%+v", err)),
	}, err)
}

// BadRequest returns a plaintext result.
func (trp TextResultProvider) BadRequest(err error) Result {
	if err != nil {
		return &RawResult{
			StatusCode:  http.StatusBadRequest,
			ContentType: ContentTypeText,
			Response:    []byte(fmt.Sprintf("Bad Request: %v", err)),
		}
	}
	return &RawResult{
		StatusCode:  http.StatusBadRequest,
		ContentType: ContentTypeText,
		Response:    []byte("Bad Request"),
	}
}

// Status returns a plaintext result.
func (trp TextResultProvider) Status(statusCode int, response ...interface{}) Result {
	return &RawResult{
		StatusCode:  statusCode,
		ContentType: ContentTypeText,
		Response:    []byte(fmt.Sprintf("%v", ResultOrDefault(http.StatusText(statusCode), response...))),
	}
}

// OK returns an plaintext result.
func (trp TextResultProvider) OK() Result {
	return &RawResult{
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeText,
		Response:    []byte("OK!"),
	}
}

// Result returns a plaintext result.
func (trp TextResultProvider) Result(result interface{}) Result {
	return &RawResult{
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeText,
		Response:    []byte(fmt.Sprintf("%v", result)),
	}
}
