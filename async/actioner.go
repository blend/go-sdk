/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "context"

// Actioner is a type that can be used as a tracked action.
type Actioner interface {
	Action(context.Context, interface{}) (interface{}, error)
}

var (
	_ Actioner = (*ActionerFunc)(nil)
)

// ActionerFunc is a function that implements action.
type ActionerFunc func(context.Context, interface{}) (interface{}, error)

// Action implements actioner for the function.
func (af ActionerFunc) Action(ctx context.Context, args interface{}) (interface{}, error) {
	return af(ctx, args)
}

var (
	_ Actioner = (*NoopActioner)(nil)
)

// NoopActioner is an actioner type that does nothing.
type NoopActioner struct{}

// Action implements actioner
func (n NoopActioner) Action(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}
