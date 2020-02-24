package configutil

import "context"

// BareResolver is the legacy / deprecated interface.
// Please add support for the context-ful resolver.
type BareResolver interface {
	Resolve() error
}

// Resolver is a type that can be resolved.
type Resolver interface {
	Resolve(context.Context) error
}
