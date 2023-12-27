/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

// RouteParameters are parameters sourced from parsing the request path (route).
type RouteParameters map[string]string

// Get gets a value for a key.
func (rp RouteParameters) Get(key string) string {
	return rp[key]
}

// Has returns if the collection has a key or not.
func (rp RouteParameters) Has(key string) bool {
	_, ok := rp[key]
	return ok
}

// Set stores a value for a key.
func (rp RouteParameters) Set(key, value string) {
	rp[key] = value
}
