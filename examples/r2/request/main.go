/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/blend/go-sdk/r2"
)

func makeRequest(path string, arguments ...r2.Option) (*http.Response, error) {
	fullOptions := append(arguments,
		r2.OptPath(path),
		r2.OptQueryValue("limit", "10"),
		r2.OptQueryValue("offset", "100"),
	)

	return r2.New("http://localhost:5000", fullOptions...).Do()
}

func main() {

	page := 10
	pageSize := 50
	opts := []r2.Option{
		r2.OptQueryValue("limit", strconv.Itoa(page)),
		r2.OptQueryValue("offset", strconv.Itoa(page*pageSize)),
	}

	res, err := makeRequest("/headers", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
	os.Exit(0)
}
