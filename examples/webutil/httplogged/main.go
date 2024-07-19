/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

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

func subScopeHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"did sub-context things"}`))
}

func scopeMetaHandler(res http.ResponseWriter, req *http.Request) {
	*req = *req.WithContext(logger.WithLabels(req.Context(), logger.Labels{"foo": "bar"}))
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status":"ok!"}`))
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

	log := logger.Prod()

	http.HandleFunc("/", webutil.HTTPLogged(log)(indexHandler))

	http.HandleFunc("/fatalerror", webutil.HTTPLogged(log)(fatalErrorHandler))
	http.HandleFunc("/error", webutil.HTTPLogged(log)(errorHandler))
	http.HandleFunc("/warning", webutil.HTTPLogged(log)(warningHandler))
	http.HandleFunc("/audit", webutil.HTTPLogged(log)(auditHandler))

	http.HandleFunc("/subscope", webutil.HTTPLogged(log.WithPath("a sub scope"))(subScopeHandler))
	http.HandleFunc("/scopemeta", webutil.HTTPLogged(log)(scopeMetaHandler))

	http.HandleFunc("/bench/logged", webutil.HTTPLogged(log)(indexHandler))

	log.Infof("Listening on :%s", port())
	log.Infof("Events %s", log.Flags.String())

	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
