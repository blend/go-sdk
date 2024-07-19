/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/blend/go-sdk/statsd"
)

var (
	flagBindAddr = flag.String("bind-addr", bindAddr(), "The bind address for the server")
)

func bindAddr() string {
	if value := os.Getenv("BIND_ADDR"); value != "" {
		return value
	}
	return "127.0.0.1:0"
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "statsd|", log.LstdFlags)

	server := &statsd.Server{
		Addr: *flagBindAddr,
		Log:  logger,
		Handler: func(metrics ...statsd.Metric) {
			printer := json.NewEncoder(os.Stdout)
			printer.SetIndent("", "  ")
			for _, metric := range metrics {
				_ = printer.Encode(metric)
			}
		},
	}
	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}
}
