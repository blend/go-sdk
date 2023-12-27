/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

// Route is an entry in the route tree.
type Route struct {
	Handler
	Method	string
	Path	string
	Params	[]string
}

// String returns the path.
func (r Route) String() string	{ return r.Path }

// StringWithMethod returns a string representation of the route.
// Namely: Method_Path
func (r Route) StringWithMethod() string {
	return r.Method + "_" + r.Path
}
