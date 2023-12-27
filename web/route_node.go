/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package web

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// RouteNodeType is a type of route node.
type RouteNodeType uint8

// RouteNodeTypes
const (
	RouteNodeTypeStatic RouteNodeType = iota // default
	RouteNodeTypeRoot
	RouteNodeTypeParam
	RouteNodeTypeCatchAll
)

// RouteNode is a node on the route tree.
type RouteNode struct {
	RouteNodeType

	Path       string
	IsWildcard bool
	MaxParams  uint8
	Indices    string
	Children   []*RouteNode
	Route      *Route
	Priority   uint32
}

// GetPath returns the node for a path, parameter values, and if there is a trailing slash redirect
// recommendation.
func (n *RouteNode) GetPath(path string) (route *Route, p RouteParameters, tsr bool) {
	return n.getValue(path)
}

// incrementChildPriority increments priority of the given child and reorders if necessary
func (n *RouteNode) incrementChildPriority(index int) int {
	n.Children[index].Priority++
	priority := n.Children[index].Priority

	// adjust position (move to front)
	newIndex := index
	for newIndex > 0 && n.Children[newIndex-1].Priority < priority {
		// swap node positions
		temp := n.Children[newIndex-1]
		n.Children[newIndex-1] = n.Children[newIndex]
		n.Children[newIndex] = temp
		newIndex--
	}

	// build new index char string
	if newIndex != index {
		n.Indices = n.Indices[:newIndex] + // unchanged prefix, might be empty
			n.Indices[index:index+1] + // the index char we move
			n.Indices[newIndex:index] + n.Indices[index+1:] // rest without char at 'pos'
	}

	return newIndex
}

// AddRoute adds a node with the given handle to the path.
func (n *RouteNode) AddRoute(method, path string, handler Handler) {
	fullPath := path
	n.Priority++
	numParams := countParams(path)

	// non-empty tree
	if len(n.Path) > 0 || len(n.Children) > 0 {
	walk:
		for {
			// Update maxParams of the current node
			if numParams > n.MaxParams {
				n.MaxParams = numParams
			}

			// Find the longest common prefix.
			// This also implies that the common prefix contains no ':' or '*'
			// since the existing key can't contain those chars.
			i := 0
			max := min(len(path), len(n.Path))
			for i < max && path[i] == n.Path[i] {
				i++
			}

			// Split edge
			if i < len(n.Path) {
				child := RouteNode{
					Path:          n.Path[i:],
					IsWildcard:    n.IsWildcard,
					RouteNodeType: RouteNodeTypeStatic,
					Indices:       n.Indices,
					Children:      n.Children,
					Route:         n.Route,
					Priority:      n.Priority - 1,
				}

				// Update maxParams (max of all Children)
				for i := range child.Children {
					if child.Children[i].MaxParams > child.MaxParams {
						child.MaxParams = child.Children[i].MaxParams
					}
				}

				n.Children = []*RouteNode{&child}
				// []byte for proper unicode char conversion, see #65
				n.Indices = string([]byte{n.Path[i]})
				n.Path = path[:i]
				n.Route = nil
				n.IsWildcard = false
			}

			// Make new node a child of this node
			if i < len(path) {
				path = path[i:]

				if n.IsWildcard {
					n = n.Children[0]
					n.Priority++

					// Update maxParams of the child node
					if numParams > n.MaxParams {
						n.MaxParams = numParams
					}
					numParams--

					// Check if the wildcard matches
					if len(path) >= len(n.Path) && n.Path == path[:len(n.Path)] {
						// check for longer wildcard, e.g. :name and :names
						if len(n.Path) >= len(path) || path[len(n.Path)] == '/' {
							continue walk
						}
					}

					panic("path segment '" + path +
						"' conflicts with existing wildcard '" + n.Path +
						"' in path '" + fullPath + "'")
				}

				c := path[0]

				// slash after param
				if n.RouteNodeType == RouteNodeTypeParam && c == '/' && len(n.Children) == 1 {
					n = n.Children[0]
					n.Priority++
					continue walk
				}

				// Check if a child with the next path byte exists
				for i := 0; i < len(n.Indices); i++ {
					if c == n.Indices[i] {
						i = n.incrementChildPriority(i)
						n = n.Children[i]
						continue walk
					}
				}

				// Otherwise insert it
				if c != ':' && c != '*' {
					// []byte for proper unicode char conversion, see #65
					n.Indices += string([]byte{c})
					child := &RouteNode{
						MaxParams: numParams,
					}
					n.Children = append(n.Children, child)
					n.incrementChildPriority(len(n.Indices) - 1)
					n = child
				}
				n.insertChild(numParams, method, path, fullPath, handler)
				return
			} else if i == len(path) { // Make node a (in-path) leaf
				if n.Route != nil {
					panic("a handle is already registered for path '" + fullPath + "'")
				}
				n.Route = &Route{
					Handler: handler,
					Path:    fullPath,
					Method:  method,
				}
			}
			return
		}
	} else { // Empty tree
		n.insertChild(numParams, method, path, fullPath, handler)
		n.RouteNodeType = RouteNodeTypeRoot
	}
}

func (n *RouteNode) insertChild(numParams uint8, method, path, fullPath string, handler Handler) {
	var offset int // already handled bytes of the path

	// find prefix until first wildcard (beginning with ':'' or '*'')
	for i, max := 0, len(path); numParams > 0; i++ {
		c := path[i]
		if c != ':' && c != '*' {
			continue
		}

		// find wildcard end (either '/' or path end)
		end := i + 1
		for end < max && path[end] != '/' {
			switch path[end] {
			// the wildcard name must not contain ':' and '*'
			case ':', '*':
				panic("only one wildcard per path segment is allowed, has: '" +
					path[i:] + "' in path '" + fullPath + "'")
			default:
				end++
			}
		}

		// check if this Node existing Children which would be
		// unreachable if we insert the wildcard here
		if len(n.Children) > 0 {
			panic("wildcard route '" + path[i:end] +
				"' conflicts with existing Children in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if end-i < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if c == ':' { // param
			// split path at the beginning of the wildcard
			if i > 0 {
				n.Path = path[offset:i]
				offset = i
			}

			child := &RouteNode{
				RouteNodeType: RouteNodeTypeParam,
				MaxParams:     numParams,
			}
			n.Children = []*RouteNode{child}
			n.IsWildcard = true
			n = child
			n.Priority++
			numParams--

			// if the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if end < max {
				n.Path = path[offset:end]
				offset = end

				child := &RouteNode{
					MaxParams: numParams,
					Priority:  1,
				}
				n.Children = []*RouteNode{child}
				n = child
			}
		} else { // catchAll
			if end != max || numParams > 1 {
				panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
			}

			if len(n.Path) > 0 && n.Path[len(n.Path)-1] == '/' {
				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
			}

			// currently fixed width 1 for '/'
			i--
			if path[i] != '/' {
				panic("no / before catch-all in path '" + fullPath + "'")
			}

			n.Path = path[offset:i]

			// first node: catchAll node with empty path
			child := &RouteNode{
				IsWildcard:    true,
				RouteNodeType: RouteNodeTypeCatchAll,
				MaxParams:     1,
			}
			n.Children = []*RouteNode{child}
			n.Indices = string(path[i])
			n = child
			n.Priority++

			// second node: node holding the variable
			child = &RouteNode{
				Path:          path[i:],
				RouteNodeType: RouteNodeTypeCatchAll,
				MaxParams:     1,
				Route: &Route{
					Handler: handler,
					Path:    fullPath,
					Method:  method,
				},
				Priority: 1,
			}
			n.Children = []*RouteNode{child}

			return
		}
	}

	// insert remaining path part and handle to the leaf
	n.Path = path[offset:]
	n.Route = &Route{
		Handler: handler,
		Path:    fullPath,
		Method:  method,
	}
}

// getValue Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *RouteNode) getValue(path string) (route *Route, p RouteParameters, tsr bool) {
walk: // outer loop for walking the tree
	for {
		if len(path) > len(n.Path) {
			if path[:len(n.Path)] == n.Path {
				path = path[len(n.Path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.IsWildcard {
					c := path[0]
					for i := 0; i < len(n.Indices); i++ {
						if c == n.Indices[i] {
							n = n.Children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = (path == "/" && n.Route != nil)
					return
				}

				// handle wildcard child
				n = n.Children[0]
				switch n.RouteNodeType {
				case RouteNodeTypeParam:
					// find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// save param value
					if p == nil {
						// lazy allocation
						p = make(RouteParameters)
					}
					p[n.Path[1:]] = path[:end]

					// we need to go deeper!
					if end < len(path) {
						if len(n.Children) > 0 {
							path = path[end:]
							n = n.Children[0]
							continue walk
						}

						// ... but we can't
						tsr = (len(path) == end+1)
						return
					}

					if route = n.Route; route != nil {
						return
					} else if len(n.Children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.Children[0]
						tsr = (n.Path == "/" && n.Route != nil)
					}

					return

				case RouteNodeTypeCatchAll:
					// save param value
					if p == nil {
						// lazy allocation
						p = make(RouteParameters)
					}

					// translation note:
					// was path[:] but the effect is the same
					p[n.Path[2:]] = path

					route = n.Route
					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == n.Path {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if route = n.Route; route != nil {
				return
			}

			if path == "/" && n.IsWildcard && n.RouteNodeType != RouteNodeTypeRoot {
				tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i := 0; i < len(n.Indices); i++ {
				if n.Indices[i] == '/' {
					n = n.Children[i]
					tsr = (len(n.Path) == 1 && n.Route != nil) ||
						(n.RouteNodeType == RouteNodeTypeCatchAll && n.Children[0].Route != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(n.Path) == len(path)+1 && n.Path[len(path)] == '/' &&
				path == n.Path[:len(n.Path)-1] && n.Route != nil)
		return
	}
}

// Makes a case-insensitive lookup of the given path and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (n *RouteNode) findCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool) {
	return n.findCaseInsensitivePathRec(
		path,
		strings.ToLower(path),
		make([]byte, 0, len(path)+1), // preallocate enough memory for new path
		[4]byte{},                    // empty rune buffer
		fixTrailingSlash,
	)
}

// shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

// recursive case-insensitive lookup function used by n.findCaseInsensitivePath
func (n *RouteNode) findCaseInsensitivePathRec(path, loPath string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) ([]byte, bool) {
	loNPath := strings.ToLower(n.Path)

walk: // outer loop for walking the tree
	for len(loPath) >= len(loNPath) && (len(loNPath) == 0 || loPath[1:len(loNPath)] == loNPath[1:]) {
		// add common path to result
		ciPath = append(ciPath, n.Path...)

		if path = path[len(n.Path):]; len(path) > 0 {
			loOld := loPath
			loPath = loPath[len(loNPath):]

			// If this node does not have a wildcard (param or catchAll) child,
			// we can just look up the next child node and continue to walk down
			// the tree
			if !n.IsWildcard {
				// skip rune bytes already processed
				rb = shiftNRuneBytes(rb, len(loNPath))

				if rb[0] != 0 {
					// old rune not finished
					for i := 0; i < len(n.Indices); i++ {
						if n.Indices[i] == rb[0] {
							// continue with child node
							n = n.Children[i]
							loNPath = strings.ToLower(n.Path)
							continue walk
						}
					}
				} else {
					// process a new rune
					var rv rune

					// find rune start
					// runes are up to 4 byte long,
					// -4 would definitely be another rune
					var off int
					for max := min(len(loNPath), 3); off < max; off++ {
						if i := len(loNPath) - off; utf8.RuneStart(loOld[i]) {
							// read rune from cached lowercase path
							rv, _ = utf8.DecodeRuneInString(loOld[i:])
							break
						}
					}

					// calculate lowercase bytes of current rune
					utf8.EncodeRune(rb[:], rv)
					// skipp already processed bytes
					rb = shiftNRuneBytes(rb, off)

					for i := 0; i < len(n.Indices); i++ {
						// lowercase matches
						if n.Indices[i] == rb[0] {
							// must use a recursive approach since both the
							// uppercase byte and the lowercase byte might exist
							// as an index
							if out, found := n.Children[i].findCaseInsensitivePathRec(
								path, loPath, ciPath, rb, fixTrailingSlash,
							); found {
								return out, true
							}
							break
						}
					}

					// same for uppercase rune, if it differs
					if up := unicode.ToUpper(rv); up != rv {
						utf8.EncodeRune(rb[:], up)
						rb = shiftNRuneBytes(rb, off)

						for i := 0; i < len(n.Indices); i++ {
							// uppercase matches
							if n.Indices[i] == rb[0] {
								// continue with child node
								n = n.Children[i]
								loNPath = strings.ToLower(n.Path)
								continue walk
							}
						}
					}
				}

				// Nothing found. We can recommend to redirect to the same URL
				// without a trailing slash if a leaf exists for that path
				return ciPath, (fixTrailingSlash && path == "/" && n.Route != nil)
			}

			n = n.Children[0]
			switch n.RouteNodeType {
			case RouteNodeTypeParam:
				// find param end (either '/' or path end)
				k := 0
				for k < len(path) && path[k] != '/' {
					k++
				}

				// add param value to case insensitive path
				ciPath = append(ciPath, path[:k]...)

				// we need to go deeper!
				if k < len(path) {
					if len(n.Children) > 0 {
						// continue with child node
						n = n.Children[0]
						loNPath = strings.ToLower(n.Path)
						loPath = loPath[k:]
						path = path[k:]
						continue
					}

					// ... but we can't
					if fixTrailingSlash && len(path) == k+1 {
						return ciPath, true
					}
					return ciPath, false
				}

				if n.Route != nil {
					return ciPath, true
				} else if fixTrailingSlash && len(n.Children) == 1 {
					// No handle found. Check if a handle for this path + a
					// trailing slash exists
					n = n.Children[0]
					if n.Path == "/" && n.Route != nil {
						return append(ciPath, '/'), true
					}
				}
				return ciPath, false

			case RouteNodeTypeCatchAll:
				return append(ciPath, path...), true

			default:
				panic("invalid node type")
			}
		} else {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if n.Route != nil {
				return ciPath, true
			}

			// No handle found.
			// Try to fix the path by adding a trailing slash
			if fixTrailingSlash {
				for i := 0; i < len(n.Indices); i++ {
					if n.Indices[i] == '/' {
						n = n.Children[i]
						if (len(n.Path) == 1 && n.Route != nil) ||
							(n.RouteNodeType == RouteNodeTypeCatchAll && n.Children[0].Route != nil) {
							return append(ciPath, '/'), true
						}
						return ciPath, false
					}
				}
			}
			return ciPath, false
		}
	}

	// Nothing found.
	// Try to fix the path by adding / removing a trailing slash
	if fixTrailingSlash {
		if path == "/" {
			return ciPath, true
		}
		if len(loPath)+1 == len(loNPath) && loNPath[len(loPath)] == '/' &&
			loPath[1:] == loNPath[1:len(loPath)] && n.Route != nil {
			return append(ciPath, n.Path...), true
		}
	}
	return ciPath, false
}

// CleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//  1. Replace multiple slashes with a single slash.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path:
//     that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
func CleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 2 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r++

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 2

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}

func countParams(path string) uint8 {
	var n uint
	for i := 0; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= 255 {
		return 255
	}
	return uint8(n)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
