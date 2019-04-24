package secrets

// KV is a basic key value store.
type KV interface {
	Put(key string, data Values, options ...Option) error
	Get(key string, options ...Option) (Values, error)
	Delete(key string, options ...Option) error
}
