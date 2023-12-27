/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func createRequestEvent(req *http.Request, rw webutil.ResponseWriter, start time.Time) webutil.HTTPRequestEvent {
	return webutil.NewHTTPRequestEvent(req,
		webutil.OptHTTPRequestStatusCode(rw.StatusCode()),
		webutil.OptHTTPRequestContentLength(rw.ContentLength()),
		webutil.OptHTTPRequestElapsed(time.Since(start)),
	)
}

func logged(log logger.Log, handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rw := webutil.NewStatusResponseWriter(res)
		handler(rw, req)
		log.TriggerContext(req.Context(), createRequestEvent(req, rw, start))
	}
}

func stdoutLogged(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler(res, req)
		fmt.Printf("%s %s %s %s %s %s %s\n",
			time.Now().UTC().Format(time.RFC3339),
			"web.request",
			req.Method,
			req.URL.Path,
			"200",
			time.Since(start).String(),
			"??",
		)
	}
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"ok!"}`))
}

func fatalErrorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(`{"status":"not ok."}`))
}

func errorHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(`{"status":"not ok."}`))
}

func warningHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{"status":"not ok."}`))
}

func subContextHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"did sub-context things"}`))
}

func auditHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"audit logged!"}`))
}

func port() string {
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		return envPort
	}
	return "8888"
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log := logger.Prod(logger.OptJSON())

	http.HandleFunc("/", logged(log, indexHandler))

	http.HandleFunc("/sub-context", logged(log, subContextHandler))
	http.HandleFunc("/fatalerror", logged(log, fatalErrorHandler))
	http.HandleFunc("/error", logged(log, errorHandler))
	http.HandleFunc("/warning", logged(log, warningHandler))
	http.HandleFunc("/audit", logged(log, auditHandler))

	http.HandleFunc("/bench/logged", logged(log, indexHandler))
	http.HandleFunc("/bench/stdout", stdoutLogged(indexHandler))

	log.Infof("Listening on :%s", port())
	log.Infof("Events %s", log.Flags.String())
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
