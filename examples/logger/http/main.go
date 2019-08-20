package main

import (
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/blend/go-sdk/logger"
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

	http.HandleFunc("/", logger.HTTPLogged(log)(indexHandler))

	http.HandleFunc("/fatalerror", logger.HTTPLogged(log)(fatalErrorHandler))
	http.HandleFunc("/error", logger.HTTPLogged(log)(errorHandler))
	http.HandleFunc("/warning", logger.HTTPLogged(log)(warningHandler))
	http.HandleFunc("/audit", logger.HTTPLogged(log)(auditHandler))

	http.HandleFunc("/subscope", logger.HTTPLogged(log.WithPath("a sub scope"))(subScopeHandler))
	http.HandleFunc("/scopemeta", logger.HTTPLogged(log)(scopeMetaHandler))

	http.HandleFunc("/bench/logged", logger.HTTPLogged(log)(indexHandler))

	log.Infof("Listening on :%s", port())
	log.Infof("Events %s", log.Flags.String())

	log.Fatal(http.ListenAndServe(":"+port(), nil))
}
