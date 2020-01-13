package main

import (
	"context"
	"flag"
	"fmt"

	"google.golang.org/grpc"

	full "github.com/blend/go-sdk/examples/grpcutil/full/protos/v1"
	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/logger"
)

var (
	addr = flag.String("addr", "localhost:9000", "the server address")
)

func main() {
	flag.Parse()

	conn, err := grpcutil.DialAddress(*addr, grpc.WithInsecure())
	if err != nil {
		logger.FatalExit(err)
	}

	client := full.NewStatusClient(conn)

	res, err := client.Status(context.Background(), &full.StatusArgs{
		MinVersion: 1,
	})
	if err != nil {
		logger.FatalExit(err)
	}
	fmt.Printf("%#v\n", *res)
}
