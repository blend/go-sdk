/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	// DefaultRetriableCodes is a set of well known types gRPC codes that should be retri-able.
	//
	// `ResourceExhausted` means that the user quota, e.g. per-RPC limits, have been reached.
	// `Unavailable` means that system is currently unavailable and the client should retry again.
	DefaultRetriableCodes	= []codes.Code{codes.ResourceExhausted, codes.Unavailable}

	defaultRetryOptions	= &retryOptions{
		max:		0,	// disabled
		perCallTimeout:	0,	// disabled
		includeHeader:	true,
		codes:		DefaultRetriableCodes,
		backoffFunc: BackoffFuncContext(func(ctx context.Context, attempt uint) time.Duration {
			return BackoffLinearWithJitter(50*time.Millisecond, 0.10)(attempt)
		}),
	}
)

// Metadata Keys
const (
	MetadataKeyAttempt = "x-retry-attempty"
)

// WithRetriesDisabled disables the retry behavior on this call, or this interceptor.
//
// Its semantically the same to `WithMax`
func WithRetriesDisabled() CallOption {
	return WithClientRetries(0)
}

// WithClientRetries sets the maximum number of retries on this call, or this interceptor.
func WithClientRetries(maxRetries uint) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.max = maxRetries
	}}
}

// WithClientRetryBackoffLinear sets the retry backoff to a fixed duration.
func WithClientRetryBackoffLinear(d time.Duration) CallOption {
	return WithClientRetryBackoffFunc(BackoffLinear(d))
}

// WithClientRetryBackoffFunc sets the `ClientRetryBackoffFunc` used to control time between retries.
func WithClientRetryBackoffFunc(bf BackoffFunc) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.backoffFunc = BackoffFuncContext(func(ctx context.Context, attempt uint) time.Duration {
			return bf(attempt)
		})
	}}
}

// WithClientRetryBackoffContext sets the `BackoffFuncContext` used to control time between retries.
func WithClientRetryBackoffContext(bf BackoffFuncContext) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.backoffFunc = bf
	}}
}

// WithClientRetryCodes sets which codes should be retried.
//
// Please *use with care*, as you may be retrying non-idempotent calls.
//
// You cannot automatically retry on Canceled and Deadline, please use `WithPerRetryTimeout` for these.
func WithClientRetryCodes(retryCodes ...codes.Code) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.codes = retryCodes
	}}
}

// WithClientRetryPerRetryTimeout sets the RPC timeout per call (including initial call) on this call, or this interceptor.
//
// The context.Deadline of the call takes precedence and sets the maximum time the whole invocation
// will take, but WithPerRetryTimeout can be used to limit the RPC time per each call.
//
// For example, with context.Deadline = now + 10s, and WithPerRetryTimeout(3 * time.Seconds), each
// of the retry calls (including the initial one) will have a deadline of now + 3s.
//
// A value of 0 disables the timeout overrides completely and returns to each retry call using the
// parent `context.Deadline`.
//
// Note that when this is enabled, any DeadlineExceeded errors that are propagated up will be retried.
func WithClientRetryPerRetryTimeout(timeout time.Duration) CallOption {
	return CallOption{applyFunc: func(o *retryOptions) {
		o.perCallTimeout = timeout
	}}
}

type retryOptions struct {
	max		uint
	perCallTimeout	time.Duration
	includeHeader	bool
	codes		[]codes.Code
	backoffFunc	BackoffFuncContext
	abortOnFailure	bool
}

// CallOption is a grpc.CallOption that is local to grpc_retry.
type CallOption struct {
	grpc.EmptyCallOption	// make sure we implement private after() and before() fields so we don't panic.
	applyFunc		func(opt *retryOptions)
}

func reuseOrNewWithCallOptions(opt *retryOptions, callOptions []CallOption) *retryOptions {
	if len(callOptions) == 0 {
		return opt
	}
	optCopy := new(retryOptions)
	*optCopy = *opt
	for _, f := range callOptions {
		f.applyFunc(optCopy)
	}
	return optCopy
}

func filterCallOptions(callOptions []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []CallOption) {
	for _, opt := range callOptions {
		if co, ok := opt.(CallOption); ok {
			retryOptions = append(retryOptions, co)
		} else {
			grpcOptions = append(grpcOptions, opt)
		}
	}
	return grpcOptions, retryOptions
}

// RetryUnaryClientInterceptor returns a new retrying unary client interceptor.
//
// The default configuration of the interceptor is to not retry *at all*. This behavior can be
// changed through options (e.g. WithMax) on creation of the interceptor or on call (through grpc.CallOptions).
func RetryUnaryClientInterceptor(optFuncs ...CallOption) grpc.UnaryClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultRetryOptions, optFuncs)
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		if callOpts.max == 0 {
			return invoker(parentCtx, method, req, reply, cc, grpcOpts...)
		}
		var lastErr error
		for attempt := uint(0); attempt < callOpts.max; attempt++ {
			callCtx, cancel := perCallContext(parentCtx, callOpts, attempt)
			func() {
				defer cancel()
				lastErr = invoker(callCtx, method, req, reply, cc, grpcOpts...)
			}()
			if lastErr == nil {
				return nil
			}
			if isContextError(lastErr) {
				if parentCtx.Err() != nil {
					// its the parent context deadline or cancellation.
					return lastErr
				} else if callOpts.perCallTimeout != 0 {
					// We have set a perCallTimeout in the retry middleware, which would result in a context error if
					// the deadline was exceeded, in which case try again.
					continue
				}
			}
			if !isRetriable(lastErr, callOpts) {
				return lastErr
			}
			if err := waitRetryBackoff(parentCtx, attempt, callOpts); err != nil {
				return err
			}
		}
		return lastErr
	}
}

// RetryStreamClientInterceptor returns a new retrying stream client interceptor for server side streaming calls.
//
// The default configuration of the interceptor is to not retry *at all*. This behavior can be
// changed through options (e.g. WithMax) on creation of the interceptor or on call (through grpc.CallOptions).
//
// Retry logic is available *only for ServerStreams*, i.e. 1:n streams, as the internal logic needs
// to buffer the messages sent by the client. If retry is enabled on any other streams (ClientStreams,
// BidiStreams), the retry interceptor will fail the call.
func RetryStreamClientInterceptor(optFuncs ...CallOption) grpc.StreamClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultRetryOptions, optFuncs)
	return func(parentCtx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		// short circuit for simplicity, and avoiding allocations.
		if callOpts.max == 0 {
			return streamer(parentCtx, desc, cc, method, grpcOpts...)
		}
		if desc.ClientStreams {
			return nil, status.Errorf(codes.Unimplemented, "grpc_retry: cannot retry on ClientStreams, set grpc_retry.Disable()")
		}

		var lastErr error
		for attempt := uint(0); attempt < callOpts.max; attempt++ {
			if err := waitRetryBackoff(parentCtx, attempt, callOpts); err != nil {
				return nil, err
			}
			callCtx, cancel := perCallContext(parentCtx, callOpts, 0)

			var newStreamer grpc.ClientStream
			func() {
				defer cancel()
				newStreamer, lastErr = streamer(callCtx, desc, cc, method, grpcOpts...)
			}()
			if lastErr == nil {
				retryingStreamer := &serverStreamingRetryingStream{
					ClientStream:	newStreamer,
					callOpts:	callOpts,
					parentCtx:	parentCtx,
					streamerCall: func(ctx context.Context) (grpc.ClientStream, error) {
						return streamer(ctx, desc, cc, method, grpcOpts...)
					},
				}
				return retryingStreamer, nil
			}

			if isContextError(lastErr) {
				if parentCtx.Err() != nil {
					// its the parent context deadline or cancellation.
					return nil, lastErr
				} else if callOpts.perCallTimeout != 0 {
					// We have set a perCallTimeout in the retry middleware, which would result in a context error if
					// the deadline was exceeded, in which case try again.
					continue
				}
			}
			if !isRetriable(lastErr, callOpts) {
				return nil, lastErr
			}
		}
		return nil, lastErr
	}
}

// type serverStreamingRetryingStream is the implementation of grpc.ClientStream that acts as a
// proxy to the underlying call. If any of the RecvMsg() calls fail, it will try to reestablish
// a new ClientStream according to the retry policy.
type serverStreamingRetryingStream struct {
	grpc.ClientStream
	bufferedSends	[]interface{}	// single message that the client can sen
	receivedGood	bool		// indicates whether any prior receives were successful
	wasClosedSend	bool		// indicates that CloseSend was closed
	parentCtx	context.Context
	callOpts	*retryOptions
	streamerCall	func(ctx context.Context) (grpc.ClientStream, error)
	mu		sync.RWMutex
}

func (s *serverStreamingRetryingStream) setStream(clientStream grpc.ClientStream) {
	s.mu.Lock()
	s.ClientStream = clientStream
	s.mu.Unlock()
}

func (s *serverStreamingRetryingStream) getStream() grpc.ClientStream {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ClientStream
}

func (s *serverStreamingRetryingStream) SendMsg(m interface{}) error {
	s.mu.Lock()
	s.bufferedSends = append(s.bufferedSends, m)
	s.mu.Unlock()
	return s.getStream().SendMsg(m)
}

func (s *serverStreamingRetryingStream) CloseSend() error {
	s.mu.Lock()
	s.wasClosedSend = true
	s.mu.Unlock()
	return s.getStream().CloseSend()
}

func (s *serverStreamingRetryingStream) Header() (metadata.MD, error) {
	return s.getStream().Header()
}

func (s *serverStreamingRetryingStream) Trailer() metadata.MD {
	return s.getStream().Trailer()
}

func (s *serverStreamingRetryingStream) RecvMsg(m interface{}) error {
	attemptRetry, lastErr := s.receiveMsgAndIndicateRetry(m)
	if !attemptRetry {
		return lastErr	// success or hard failure
	}
	// We start off from attempt 1, because zeroth was already made on normal SendMsg().
	for attempt := uint(1); attempt < s.callOpts.max; attempt++ {
		if err := waitRetryBackoff(s.parentCtx, attempt, s.callOpts); err != nil {
			return err
		}
		callCtx, cancel := perCallContext(s.parentCtx, s.callOpts, attempt)

		var newStream grpc.ClientStream
		var err error
		func() {
			defer cancel()
			newStream, err = s.reestablishStreamAndResendBuffer(callCtx)
		}()
		if err != nil {
			// TODO(mwitkow): Maybe dial and transport errors should be retriable?
			return err
		}
		s.setStream(newStream)
		attemptRetry, lastErr = s.receiveMsgAndIndicateRetry(m)
		//fmt.Printf("Received message and indicate: %v  %v\n", attemptRetry, lastErr)
		if !attemptRetry {
			return lastErr
		}
	}
	return lastErr
}

func (s *serverStreamingRetryingStream) receiveMsgAndIndicateRetry(m interface{}) (bool, error) {
	s.mu.RLock()
	wasGood := s.receivedGood
	s.mu.RUnlock()
	err := s.getStream().RecvMsg(m)
	if err == nil || err == io.EOF {
		s.mu.Lock()
		s.receivedGood = true
		s.mu.Unlock()
		return false, err
	} else if wasGood {
		// previous RecvMsg in the stream succeeded, no retry logic should interfere
		return false, err
	}
	if isContextError(err) {
		if s.parentCtx.Err() != nil {
			return false, err
		} else if s.callOpts.perCallTimeout != 0 {
			// We have set a perCallTimeout in the retry middleware, which would result in a context error if
			// the deadline was exceeded, in which case try again.
			return true, err
		}
	}
	return isRetriable(err, s.callOpts), err
}

func (s *serverStreamingRetryingStream) reestablishStreamAndResendBuffer(callCtx context.Context) (grpc.ClientStream, error) {
	s.mu.RLock()
	bufferedSends := s.bufferedSends
	s.mu.RUnlock()
	newStream, err := s.streamerCall(callCtx)
	if err != nil {
		return nil, err
	}
	for _, msg := range bufferedSends {
		if err := newStream.SendMsg(msg); err != nil {
			return nil, err
		}
	}
	if err := newStream.CloseSend(); err != nil {
		return nil, err
	}
	return newStream, nil
}

func waitRetryBackoff(parentCtx context.Context, attempt uint, callOpts *retryOptions) error {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = callOpts.backoffFunc(parentCtx, attempt)
	}
	if waitTime > 0 {
		timer := time.NewTimer(waitTime)
		select {
		case <-parentCtx.Done():
			timer.Stop()
			return contextErrToGrpcErr(parentCtx.Err())
		case <-timer.C:
		}
	}
	return nil
}

func isRetriable(err error, callOpts *retryOptions) bool {
	if isContextError(err) {
		return false
	}

	errCode := status.Code(err)
	for _, code := range callOpts.codes {
		if code == errCode {
			return true
		}
	}
	return !callOpts.abortOnFailure
}

func isContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}

func perCallContext(parentCtx context.Context, callOpts *retryOptions, attempt uint) (ctx context.Context, cancel func()) {
	ctx = parentCtx
	cancel = func() {}
	if callOpts.perCallTimeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, callOpts.perCallTimeout)
	}
	if attempt > 0 && callOpts.includeHeader {
		mdClone := cloneMetadata(extractOutgoingMetadata(ctx))
		mdClone = setMetadata(mdClone, MetadataKeyAttempt, fmt.Sprintf("%d", attempt))
		ctx = toOutgoing(ctx, mdClone)
	}
	return
}

func contextErrToGrpcErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Errorf(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Errorf(codes.Canceled, err.Error())
	default:
		return status.Errorf(codes.Unknown, err.Error())
	}
}

// extractOutgoingMetadata extracts an outbound metadata from the client-side context.
//
// This function always returns a NiceMD wrapper of the metadata.MD, in case the context doesn't have metadata it returns
// a new empty NiceMD.
func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.Pairs()	// empty md set
	}
	return md
}

// cloneMetadata clones a given md set.
func cloneMetadata(md metadata.MD, copiedKeys ...string) metadata.MD {
	newMd := make(metadata.MD)
	for k, vv := range md {
		var found bool
		if len(copiedKeys) == 0 {
			found = true
		} else {
			for _, allowedKey := range copiedKeys {
				if strings.EqualFold(allowedKey, k) {
					found = true
					break
				}
			}
		}
		if !found {
			continue
		}
		newMd[k] = make([]string, len(vv))
		copy(newMd[k], vv)
	}
	return newMd
}

func setMetadata(md metadata.MD, key string, value string) metadata.MD {
	k, v := encodeMetadataKeyValue(key, value)
	md[k] = []string{v}
	return md
}

// toOutgoing sets the given NiceMD as a client-side context for dispatching.
func toOutgoing(ctx context.Context, md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}

const (
	binHdrSuffix = "-bin"
)

func encodeMetadataKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, binHdrSuffix) {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}
