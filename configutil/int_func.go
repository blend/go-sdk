package configutil

import "context"

var (
	_ IntSource = (*IntFunc)(nil)
)

// IntFunc is an int value source from a commandline flag.
type IntFunc func(context.Context) (*int, error)

// Int returns an invocation of the function.
func (vf IntFunc) Int(ctx context.Context) (*int, error) {
	return vf(ctx)
}
