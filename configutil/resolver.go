package configutil

import "context"

// Resolver is a type that can be resolved.
type Resolver interface {
	Resolve(context.Context) error
}
