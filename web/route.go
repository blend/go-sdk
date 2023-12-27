/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

// Route is an entry in the route tree.
type Route struct {
	Handler
	Method string
	Path   string
	Params []string
}

// String returns the path.
func (r Route) String() string { return r.Path }

// StringWithMethod returns a string representation of the route.
// Namely: Method_Path
func (r Route) StringWithMethod() string {
	return r.Method + "_" + r.Path
}
