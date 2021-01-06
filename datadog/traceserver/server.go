package traceserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	msgp "github.com/tinylib/msgp/msgp"
)

// Server is a server for handling traces.
type Server struct {
	Addr     string
	Log      *log.Logger
	Listener net.Listener
	Server   *http.Server
	Handler  func(context.Context, ...*Span)
}

// Start starts the server.
func (ts *Server) Start() error {
	var err error
	if ts.Handler == nil {
		return fmt.Errorf("server cannot start; no handler provided")
	}
	if ts.Listener == nil && ts.Addr != "" {
		ts.Listener, err = net.Listen("tcp", ts.Addr)
		if err != nil {
			return err
		}
	}
	if ts.Listener == nil {
		return fmt.Errorf("server cannot start; no listener or addr provided")
	}

	ts.logf("trace server listening: %s", ts.Listener.Addr().String())
	ts.Server = &http.Server{
		Handler: ts,
	}
	return ts.Server.Serve(ts.Listener)
}

// Stop stops the trace server.
func (ts *Server) Stop() error {
	if ts.Server == nil {
		return nil
	}
	if err := ts.Server.Shutdown(context.Background()); err != nil {
		return err
	}
	ts.Server = nil
	return nil
}

func (ts *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	srw := &ResponseWriter{ResponseWriter: rw}

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		ts.logf("%s %s %d %v %s", req.Method, req.URL.String(), srw.StatusCode, elapsed, FormatContentLength(srw.ContentLength))
	}()

	switch req.Method {
	case http.MethodGet:
		switch req.URL.Path {
		case "/":
			ts.handleGetIndex(srw, req)
			return
		default:
		}
	case http.MethodPost:
		switch req.URL.Path {
		case "/v0.4/traces":
			ts.handlePostTraces(srw, req)
			return
		default:
		}
	default:
	}
	http.NotFound(srw, req)
}

//
// handlers
//

func (ts *Server) handleGetIndex(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Datadog Trace Echo")
}

func (ts *Server) handlePostTraces(rw http.ResponseWriter, req *http.Request) {
	var payload SpanLists
	if err := msgp.Decode(req.Body, &payload); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	for _, spanList := range payload {
		ts.Handler(req.Context(), spanList...)
	}
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "OK!")
}

//
// logging
//

func (ts *Server) logf(format string, args ...interface{}) {
	if ts.Log != nil {
		format = strings.TrimSpace(format)
		ts.Log.Printf(format+"\n", args...)
	}
}

func (ts *Server) logln(args ...interface{}) {
	if ts.Log != nil {
		ts.Log.Println(args...)
	}
}

// FormatContentLength returns a string representation of a file size in bytes.
func FormatContentLength(sizeBytes int) string {
	if sizeBytes >= 1<<30 {
		return fmt.Sprintf("%dgB", sizeBytes/(1<<30))
	} else if sizeBytes >= 1<<20 {
		return fmt.Sprintf("%dmB", sizeBytes/(1<<20))
	} else if sizeBytes >= 1<<10 {
		return fmt.Sprintf("%dkB", sizeBytes/(1<<10))
	}
	return fmt.Sprintf("%dB", sizeBytes)
}
