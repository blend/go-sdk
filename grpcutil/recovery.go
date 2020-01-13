package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type recoveryOptions struct {
	recoveryHandlerFunc RecoveryHandlerFunc
}

// RecoveryOption is a type that provides a recovery option.
type RecoveryOption func(*recoveryOptions)

// WithRecoveryHandler customizes the function for recovering from a panic.
func WithRecoveryHandler(f RecoveryHandlerFunc) RecoveryOption {
	return func(o *recoveryOptions) {
		o.recoveryHandlerFunc = f
	}
}

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (err error)

// RecoverUnary returns a new unary server interceptor for panic recovery.
func RecoverUnary(opts ...RecoveryOption) grpc.UnaryServerInterceptor {
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

// RecoverStream returns a new streaming server interceptor for panic recovery.
func RecoverStream(opts ...RecoveryOption) grpc.StreamServerInterceptor {
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
		return grpc.Errorf(codes.Internal, "%s", p)
	}
	return r(p)
}

var (
	defaultOptions = &recoveryOptions{
		recoveryHandlerFunc: nil,
	}
)

func evaluateOptions(opts []RecoveryOption) *recoveryOptions {
	optCopy := &recoveryOptions{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}
