package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/blend/go-sdk/grpcutil"
	full "github.com/blend/go-sdk/grpcutil/examples/full/protos"
	"github.com/blend/go-sdk/logger"
	"google.golang.org/grpc"
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

	res, err := client.Status(context.Background(), &full.StatusArgs{})
	if err != nil {
		logger.FatalExit(err)
	}
	fmt.Printf("%#v\n", *res)
}
