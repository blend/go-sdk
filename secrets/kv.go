package secrets

import "context"

// KV is a basic key value store.
type KV interface {
	Put(ctx context.Context, key string, data Values, options ...Option) error
	Get(ctx context.Context, key string, options ...Option) (Values, error)
	Delete(ctx context.Context, key string, options ...Option) error
	List(ctx context.Context, path string, options ...Option) ([]string, error)
}

// KVClient is a basic key value store client.
type KVClient = KV
