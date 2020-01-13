//go:generate protoc -I ./protos/v1 --go_out=plugins=grpc:./protos/v1 ./protos/v1/full.proto
//go:generate protoc -I ./protos/v2 --go_out=plugins=grpc:./protos/v2 ./protos/v2/full.proto

// Package main implements a server for the Status service and implements a number of extra features like logging and recovery.
package main

import (
	"context"
	"fmt"
	"os"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/configutil"
	full "github.com/blend/go-sdk/examples/grpcutil/full/protos/v2"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/grpcutil"
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
	full.RegisterStatusServer(server, statusServer{
		Log: log.WithPath("status"),
	})
	if err := graceful.Shutdown(grpcutil.NewGraceful(listener, server).WithLogger(log)); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// Set by LDFLAGS
var (
	Version = 1
)

var (
	_ (full.StatusServer) = (*statusServer)(nil)
)

type statusServer struct {
	Log logger.Log
}

func (ss statusServer) Status(ctx context.Context, args *full.VersionedStatusArgs) (*full.StatusResponse, error) {
	if args.MinVersion > 1 {
		return nil, grpc.Errorf(400, "bad version: %d", args.MinVersion)
	}
	logger.MaybeInfof(ss.Log, "responding with version: %d", Version)
	return &full.StatusResponse{
		Version: int64(Version),
		GitRef:  env.Env().String("CURRENT_REF", "123456"),
	}, nil
}
