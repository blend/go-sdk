//go:generate protoc -I ./protos --go_out=plugins=grpc:./protos ./protos/full.proto

// Package main implements a server for the Status service and implements a number of extra features like logging and recovery.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blend/go-sdk/ref"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/grpcutil"
	full "github.com/blend/go-sdk/examples/grpcutil/full/protos"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/proxyprotocol"
	"google.golang.org/grpc"
)

func customUnary(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("Did custom action!")
	return handler(ctx, args)
}

type config struct {
	BindAddr         string        `json:"bindAddr" yaml:"bindAddr"`
	UseProxyProtocol *bool         `json:"useProxyProtocol" yaml:"useProxyProtocol"`
	Logger           logger.Config `json:"logger" yaml:"logger"`
}

func (c *config) Resolve() error {
	return configutil.AnyError(
		configutil.SetString(&c.BindAddr, configutil.Env("BIND_ADDR"), configutil.String(c.BindAddr), configutil.String(":9000")),
		configutil.SetBool(&c.UseProxyProtocol, configutil.Env("PROXY_PROTOCOL"), configutil.Bool(c.UseProxyProtocol), configutil.Bool(ref.Bool(false))),
	)
}

func main() {
	var cfg config
	if _, err := configutil.Read(&cfg); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}
	log := logger.MustNew(logger.OptConfig(cfg.Logger))

	log.Infof("using bind address: %s", cfg.BindAddr)
	listener, err := grpcutil.Listener(cfg.BindAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// grpc interceptors
	interceptors := []grpc.UnaryServerInterceptor{
		customUnary,
		grpcutil.RecoverUnary(),
		grpcutil.LoggedUnary(log),
	}

	// start the grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcutil.UnaryChain(interceptors...)),
	}

	if cfg.UseProxyProtocol != nil && *cfg.UseProxyProtocol {
		listener = &proxyprotocol.Listener{Listener: listener}
	}

	server := grpc.NewServer(opts...)
	full.RegisterStatusServer(server, statusServer{})
	if err := graceful.Shutdown(grpcutil.NewGraceful(listener, server).WithLogger(log)); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

var (
	_ (full.StatusServer) = (*statusServer)(nil)
)

type statusServer struct{}

func (ss statusServer) Status(context.Context, *full.StatusArgs) (*full.StatusResponse, error) {
	return &full.StatusResponse{
		Version: os.Getenv("VERSION"),
		GitRef:  os.Getenv("CURRENT_REF"),
	}, nil
}
