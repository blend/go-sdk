package secrets

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
)

// MustURL creates a new url and panics on error.
func MustURL(format string, args ...interface{}) *url.URL {
	output, err := url.ParseRequestURI(fmt.Sprintf(format, args...))
	if err != nil {
		panic(err)
	}
	return output
}

// ServiceConfigPath returns the service config path from environment.
func ServiceConfigPath(serviceEnv, serviceName string) string {
	return filepath.Join("secret", "data", serviceEnv, "service", serviceName, "config")
}

// ServiceConfigPathFromEnv returns the service config path from environment.
func ServiceConfigPathFromEnv() string {
	return ServiceConfigPath(env.Env().ServiceEnv(), env.Env().ServiceName())
}

// ExceptionClassForStatus returns the exception class for a given remote status code.
func ExceptionClassForStatus(statusCode int) ex.Class {
	switch statusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden, http.StatusUnauthorized:
		return ErrUnauthorized
	default:
		return ErrServerError
	}
}
