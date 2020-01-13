package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func main() {
	//create external transport reference
	transport := &http.Transport{}

	// pass to the request
	req := r2.New("https://google.com/robots.txt", r2.OptTransport(transport))

	var res *http.Response
	var err error
	for x := 0; x < 10; x++ {
		res, err = req.Discard()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			fmt.Fprintf(os.Stdout, "%v %v %v\n", res.StatusCode, res.Status, res.ContentLength)
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Done")
}
