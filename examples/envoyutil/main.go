package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blend/go-sdk/envoyutil"
)

func bindAddr() string {
	if val := os.Getenv("BIND_ADDR"); val != "" {
		return val
	}
	return "localhost:8080"
}

func extractIdentity(xfcc envoyutil.XFCCElement) (clientIdentity string, err error) {
	clientIdentity = xfcc.By
	return
}

/*

You should test this with:

> curl localhost:8080 -H "X-Forwarded-Client-Cert:By=spiffe://cluster.local/ns/blent/sa/echo;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject=10;URI=spiffe://cluster.local/ns/blent/sa/beep"
spiffe://cluster.local/ns/blent/sa/echo

*/

func main() {
	// GET /
	http.DefaultServeMux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		clientIdentity, err := envoyutil.ExtractAndVerifyClientIdentity(r, extractIdentity)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintln(rw, clientIdentity)
	}))

	log.Printf("listening on %s\n", bindAddr())
	// start the server
	if err := http.ListenAndServe(bindAddr(), nil); err != nil {
		log.Fatal(err)
	}
}
