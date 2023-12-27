/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

// ResultProvider is the provider interface for results.
type ResultProvider interface {
	InternalError(err error) Result
	BadRequest(err error) Result
	NotFound() Result
	NotAuthorized() Result
	Status(int, interface{}) Result
}

// ResultOrDefault returns a result or a default.
func ResultOrDefault(result, defaultResult interface{}) interface{} {
	if result != nil {
		return result
	}
	return defaultResult
}
