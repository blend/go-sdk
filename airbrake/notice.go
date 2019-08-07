package airbrake

import (
	"fmt"
	"net/http"

	"github.com/airbrake/gobrake"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/webutil"
)

// NewNotice returns a new gobrake notice.
func NewNotice(err interface{}, req *http.Request) *gobrake.Notice {
	var notice *gobrake.Notice

	if exErr := ex.As(err); exErr != nil {
		var errors []gobrake.Error
		errors = append(errors, gobrake.Error{
			Type:      fmt.Sprintf("%v", ex.ErrClass(exErr)),
			Message:   exErr.Message,
			Backtrace: frames(exErr.StackTrace),
		})

		for inner := ex.As(exErr.Inner); inner != nil; inner = ex.As(inner.Inner) {
			errors = append(errors, gobrake.Error{
				Type:      fmt.Sprintf("%v", ex.ErrClass(inner)),
				Message:   fmt.Sprintf("%+v", inner),
				Backtrace: frames(inner.StackTrace),
			})
		}
		notice = &gobrake.Notice{
			Errors:  errors,
			Context: make(map[string]interface{}),
			Env:     make(map[string]interface{}),
			Session: make(map[string]interface{}),
			Params:  make(map[string]interface{}),
		}
	} else {
		notice = &gobrake.Notice{
			Errors: []gobrake.Error{{
				Type:    fmt.Sprintf("%v", err),
				Message: fmt.Sprintf("%v", err),
			}},
			Context: make(map[string]interface{}),
			Env:     make(map[string]interface{}),
			Session: make(map[string]interface{}),
			Params:  make(map[string]interface{}),
		}
	}
	for k, v := range getDefaultContext() {
		notice.Context[k] = v
	}
	notice.Context["severity"] = "error"

	// set requests minus headers
	if req != nil {
		notice.Context["url"] = req.URL.String()
		notice.Context["httpMethod"] = req.Method
		if ua := webutil.GetUserAgent(req); ua != "" {
			notice.Context["userAgent"] = ua
		}
		notice.Context["userAddr"] = webutil.GetRemoteAddr(req)
	}
	return notice
}

func frames(stack ex.StackTrace) (output []gobrake.StackFrame) {
	if typed, ok := stack.(ex.StackPointers); ok {
		var frame ex.Frame
		for _, ptr := range typed {
			frame = ex.Frame(ptr)
			output = append(output, gobrake.StackFrame{
				File: frame.File(),
				Func: frame.Func(),
				Line: frame.Line(),
			})
		}
	}
	return
}
