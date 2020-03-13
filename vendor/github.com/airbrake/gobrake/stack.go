package gobrake

import (
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// getBacktrace returns the stacktrace associated with e. If e is an
// error from the errors package its stacktrace is extracted, otherwise
// the current stacktrace is collected end returned.
func getBacktrace(e interface{}, skip int) (string, []StackFrame) {
	if err, ok := e.(stackTracer); ok {
		return backtraceFromErrorWithStackTrace(err)
	}

	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	ff := runtime.CallersFrames(pcs[:n])

	var firstPkg string
	frames := make([]StackFrame, 0)
	for {
		f, ok := ff.Next()
		if !ok {
			break
		}

		pkg, fn := splitPackageFuncName(f.Function)
		if firstPkg == "" && pkg != "runtime" {
			firstPkg = pkg
		}

		if stackFilter(pkg, fn, f.File, f.Line) {
			frames = frames[:0]
			continue
		}

		frames = append(frames, StackFrame{
			File: f.File,
			Line: f.Line,
			Func: fn,
		})
	}

	return firstPkg, frames
}

func splitPackageFuncName(funcName string) (string, string) {
	var packageName string
	if ind := strings.LastIndex(funcName, "/"); ind > 0 {
		packageName += funcName[:ind+1]
		funcName = funcName[ind+1:]
	}
	if ind := strings.Index(funcName, "."); ind > 0 {
		packageName += funcName[:ind]
		funcName = funcName[ind+1:]
	}
	return packageName, funcName
}

func stackFilter(packageName, funcName string, file string, line int) bool {
	return packageName == "runtime" && funcName == "panic"
}

// stackTraces returns the stackTrace of an error.
// It is part of the errors package public interface.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// backtraceFromErrorWithStackTrace extracts the stacktrace from e.
func backtraceFromErrorWithStackTrace(e stackTracer) (string, []StackFrame) {
	stackTrace := e.StackTrace()
	var pcs []uintptr
	for _, f := range stackTrace {
		pcs = append(pcs, uintptr(f))
	}

	ff := runtime.CallersFrames(pcs)
	var firstPkg string
	frames := make([]StackFrame, 0)
	for {
		f, ok := ff.Next()
		if !ok {
			break
		}

		pkg, fn := splitPackageFuncName(f.Function)
		if firstPkg == "" {
			firstPkg = pkg
		}

		frames = append(frames, StackFrame{
			File: f.File,
			Line: f.Line,
			Func: fn,
		})
	}

	return firstPkg, frames
}
