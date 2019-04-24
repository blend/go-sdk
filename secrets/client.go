package secrets

// Client is the general interface for a Secrets client
type Client interface {
	Put(key string, data Values, options ...Option) error
	Get(key string, options ...Option) (Values, error)
	Delete(key string, options ...Option) error
}
