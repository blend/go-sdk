package retry

import "context"

// Action is a function you can retry.
type Action func(ctx context.Context) (interface{}, error)
