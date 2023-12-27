/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"strings"

	"github.com/blend/go-sdk/ex"
)

// Error Constants
const (
	ErrMissingRouteParameters ex.Class = "missing route parameter in params"
)

// ReplacePathParameters will replace path parameters in a URL path with values
// from the passed in `params` map. Path parameters in the format of `:<param_name>`.
// Example usage: `ReplacePathParameters("/resource/:resource_id", map[string]string{"resource_id": "1234"})`
func ReplacePathParameters(str string, params map[string]string) (string, error) {
	if params == nil {
		params = make(map[string]string)
	}

	parts := strings.Split(str, "/")
	for i := range parts {
		if !strings.HasPrefix(parts[i], ":") {
			continue
		}
		pathValue, ok := params[strings.TrimPrefix(parts[i], ":")]
		if !ok {
			pathValue, ok = params[parts[i]]
			if !ok {
				return "", ex.New(ErrMissingRouteParameters, ex.OptMessagef("Missing %s", parts[i]))
			}
		}

		parts[i] = pathValue
	}

	return strings.Join(parts, "/"), nil
}
