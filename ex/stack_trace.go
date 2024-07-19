/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// StackTraceProvider is a type that can return an exception class.
type StackTraceProvider interface {
	StackTrace() StackTrace
}

// GetStackTrace is a utility method to get the current stack trace at call time.
func GetStackTrace() string {
	return fmt.Sprintf("%+v", Callers(DefaultStartDepth))
}

// Callers returns stack pointers.
func Callers(startDepth int) StackPointers {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(startDepth, pcs[:])
	var st StackPointers = pcs[0:n]
	return st
}

// StackTrace is a stack trace provider.
type StackTrace interface {
	fmt.Formatter
	Strings() []string
	String() string
}

// StackPointers is stack of uintptr stack frames from innermost (newest) to outermost (oldest).
type StackPointers []uintptr

// Format formats the stack trace.
func (st StackPointers) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				fmt.Fprintf(s, "\n%+v", Frame(f))
			}
		case s.Flag('#'):
			for _, f := range st {
				fmt.Fprintf(s, "\n%#v", Frame(f))
			}
		default:
			for _, f := range st {
				fmt.Fprintf(s, "\n%v", Frame(f))
			}
		}
	case 's':
		for _, f := range st {
			fmt.Fprintf(s, "\n%s", Frame(f))
		}
	}
}

// Strings dereferences the StackTrace as a string slice
func (st StackPointers) Strings() []string {
	res := make([]string, len(st))
	for i, frame := range st {
		res[i] = fmt.Sprintf("%+v", Frame(frame))
	}
	return res
}

// String returns a single string representation of the stack pointers.
func (st StackPointers) String() string {
	return fmt.Sprintf("%+v", st)
}

// MarshalJSON is a custom json marshaler.
func (st StackPointers) MarshalJSON() ([]byte, error) {
	return json.Marshal(st.Strings())
}

// StackStrings represents a stack trace as string literals.
type StackStrings []string

// Format formats the stack trace.
func (ss StackStrings) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range ss {
				fmt.Fprintf(s, "\n%+v", f)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []string(ss))
		default:
			for _, f := range ss {
				fmt.Fprintf(s, "\n%v", f)
			}
		}
	case 's':
		for _, f := range ss {
			fmt.Fprintf(s, "\n%v", f)
		}
	}
}

// Strings returns the stack strings as a string slice.
func (ss StackStrings) Strings() []string {
	return []string(ss)
}

// String returns a single string representation of the stack pointers.
func (ss StackStrings) String() string {
	return fmt.Sprintf("%+v", ss)
}

// MarshalJSON is a custom json marshaler.
func (ss StackStrings) MarshalJSON() ([]byte, error) {
	return json.Marshal(ss)
}

// Frame represents a program counter inside a stack frame.
type Frame uintptr

// PC returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) PC() uintptr { return uintptr(f) - 1 }

// File returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) File() string {
	fn := runtime.FuncForPC(f.PC())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.PC())
	return file
}

// Line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.PC())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.PC())
	return line
}

// Func returns the func name.
func (f Frame) Func() string {
	name := runtime.FuncForPC(f.PC()).Name()
	return funcname(name)
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source file
//	%d    source line
//	%n    function name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   path of source file relative to the compile time GOPATH
//	%+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			pc := f.PC()
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				fmt.Fprint(s, "unknown")
			} else {
				file, _ := fn.FileLine(pc)
				fname := fn.Name()
				fmt.Fprintf(s, "%s\n\t%s", fname, trimGOPATH(fname, file))
			}
		default:
			fmt.Fprint(s, path.Base(f.File()))
		}
	case 'd':
		fmt.Fprintf(s, "%d", f.Line())
	case 'n':
		name := runtime.FuncForPC(f.PC()).Name()
		fmt.Fprint(s, funcname(name))
	case 'v':
		f.Format(s, 's')
		fmt.Fprint(s, ":")
		f.Format(s, 'd')
	}
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func trimGOPATH(name, file string) string {
	// Here we want to get the source file path relative to the compile time
	// GOPATH. As of Go 1.6.x there is no direct way to know the compiled
	// GOPATH at runtime, but we can infer the number of path segments in the
	// GOPATH. We note that fn.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/sub.Type.Method
	//
	// We want to produce:
	//
	//    pkg/sub/file.go
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired output. We count separators from the end of the file
	// path until it finds two more than in the function name and then move
	// one character forward to preserve the initial path segment without a
	// leading separator.
	const sep = "/"
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			// not enough separators found, set i so that the slice expression
			// below leaves file unmodified
			i = -len(sep)
			break
		}
	}
	// get back to 0 or trim the leading separator
	file = file[i+len(sep):]
	return file
}
