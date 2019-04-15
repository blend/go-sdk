package secrets

import "context"

// KV is a basic key value store.
type KV interface {
	Put(ctx context.Context, key string, data Values, options ...RequestOption) error
	Get(ctx context.Context, key string, options ...RequestOption) (Values, error)
	Delete(ctx context.Context, key string, options ...RequestOption) error
	List(ctx context.Context, path string, options ...RequestOption) ([]string, error)
}

// KVClient is a basic key value store client.
type KVClient = KV
