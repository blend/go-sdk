/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/blend/go-sdk/datadog/traceserver"
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

	logger := log.New(os.Stdout, "traces|", log.LstdFlags)

	server := traceserver.Server{
		Addr:	*flagBindAddr,
		Log:	logger,
		Handler: func(_ context.Context, spans ...*traceserver.Span) {
			printer := json.NewEncoder(os.Stdout)
			printer.SetIndent("", "  ")
			for _, span := range spans {
				_ = printer.Encode(span)
			}
		},
	}

	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}
}
