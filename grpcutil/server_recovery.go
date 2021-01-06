package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// LoggedRecoveryHandler is a recovery handler shim.
func LoggedRecoveryHandler(log logger.Log) RecoveryHandlerFunc {
	return func(p interface{}) error {
		logger.MaybeError(log, ex.New(p))
		return status.Errorf(codes.Internal, "%+v", p)
	}
}

type serverRecoveryOptions struct {
	recoveryHandlerFunc RecoveryHandlerFunc
}

// ServerRecoveryOption is a type that provides a recovery option.
type ServerRecoveryOption func(*serverRecoveryOptions)

// WithServerRecoveryHandler customizes the function for recovering from a panic.
func WithServerRecoveryHandler(f RecoveryHandlerFunc) ServerRecoveryOption {
	return func(o *serverRecoveryOptions) {
		o.recoveryHandlerFunc = f
	}
}

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (err error)

// RecoverServerUnary returns a new unary server interceptor for panic recovery.
func RecoverServerUnary(opts ...ServerRecoveryOption) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(r, o.recoveryHandlerFunc)
			}
		}()
		return handler(ctx, req)
	}
}

// RecoverServerStream returns a new streaming server interceptor for panic recovery.
func RecoverServerStream(opts ...ServerRecoveryOption) grpc.StreamServerInterceptor {
	o := evaluateOptions(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(r, o.recoveryHandlerFunc)
			}
		}()

		return handler(srv, stream)
	}
}

func recoverFrom(p interface{}, r RecoveryHandlerFunc) error {
	if r == nil {
		return status.Errorf(codes.Internal, "%s", p)
	}
	return r(p)
}

var (
	defaultOptions = &serverRecoveryOptions{
		recoveryHandlerFunc: nil,
	}
)

func evaluateOptions(opts []ServerRecoveryOption) *serverRecoveryOptions {
	optCopy := new(serverRecoveryOptions)
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}
