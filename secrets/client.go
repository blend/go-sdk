package secrets

import (
	"context"
)

// Client is the general interface for a Secrets client
type Client interface {
	Put(ctx context.Context, key string, data Values, options ...Option) error
	Get(ctx context.Context, key string, options ...Option) (Values, error)
	Delete(ctx context.Context, key string, options ...Option) error
	List(ctx context.Context, path string, options ...Option) ([]string, error)
}
